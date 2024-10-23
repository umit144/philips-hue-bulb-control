package light_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/umit144/philips-hue-bulb-control/internal/light"
)

func setupTestServer() (*httptest.Server, *light.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			fmt.Printf("Request body: %s\n", string(body))
		}

		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/lights" {
			w.Write([]byte(`{
				"1": {
					"state": {
						"on": true,
						"bri": 254,
						"hue": 4444,
						"sat": 254
					},
					"name": "Living Room"
				},
				"2": {
					"state": {
						"on": false,
						"bri": 254,
						"hue": 4444,
						"sat": 254
					},
					"name": "Kitchen"
				}
			}`))
			return
		}

		if matches := strings.Contains(r.URL.Path, "/lights/") && strings.Contains(r.URL.Path, "/state"); matches {
			if r.Method != http.MethodPut {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			lightID := strings.Split(r.URL.Path, "/lights/")[1]
			lightID = strings.Split(lightID, "/state")[0]

			if lightID != "1" && lightID != "2" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`[{"error": {"type": 3, "address": "/lights/3", "description": "resource not found"}}]`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"success":{"/lights/1/state/on": true}}]`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`[{"error": {"type": 3, "address": "/lights/3", "description": "resource not found"}}]`))
	}))

	client := light.NewClient("test-bridge", "testuser")
	client.SetBaseURL(server.URL)
	client.SetHTTPClient(server.Client())

	return server, client
}

func TestGetAll(t *testing.T) {
	server, client := setupTestServer()
	defer server.Close()

	lights, err := client.GetAll()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(lights) != 2 {
		t.Errorf("Expected 2 lights, got %d", len(lights))
	}

	l1, exists := lights["1"]
	if !exists {
		t.Fatal("Expected light 1 to exist")
	}
	if l1.Name != "Living Room" {
		t.Errorf("Expected light 1 name to be 'Living Room', got '%s'", l1.Name)
	}
	if !l1.State.On {
		t.Error("Expected light 1 to be on")
	}

	l2, exists := lights["2"]
	if !exists {
		t.Fatal("Expected light 2 to exist")
	}
	if l2.Name != "Kitchen" {
		t.Errorf("Expected light 2 name to be 'Kitchen', got '%s'", l2.Name)
	}
	if l2.State.On {
		t.Error("Expected light 2 to be off")
	}
}

func TestToggle(t *testing.T) {
	tests := []struct {
		name    string
		lightID string
		state   bool
		wantErr bool
	}{
		{
			name:    "toggle existing light on",
			lightID: "1",
			state:   true,
			wantErr: false,
		},
		{
			name:    "toggle existing light off",
			lightID: "2",
			state:   false,
			wantErr: false,
		},
		{
			name:    "toggle non-existing light",
			lightID: "999",
			state:   true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer()
			defer server.Close()

			err := client.Toggle(tt.lightID, tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("Toggle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
