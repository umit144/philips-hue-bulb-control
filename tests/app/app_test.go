package app_test

import (
	"testing"

	"github.com/umit144/philips-hue-bulb-control/internal/app"
	"github.com/umit144/philips-hue-bulb-control/internal/config"
	"github.com/umit144/philips-hue-bulb-control/tests/mocks"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		BridgeIP: "test-bridge",
		Username: "test-user",
	}

	application := app.New(cfg)
	if application == nil {
		t.Fatal("Expected app instance, got nil")
	}
}

func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		name    string
		lightID string
		state   string
		wantErr bool
	}{
		{
			name:    "valid light and state",
			lightID: "1",
			state:   "on",
			wantErr: false,
		},
		{
			name:    "invalid light ID",
			lightID: "999",
			state:   "on",
			wantErr: true,
		},
		{
			name:    "valid light with valid state off",
			lightID: "2",
			state:   "off",
			wantErr: false,
		},
	}

	cfg := &config.Config{
		BridgeIP: "test-bridge",
		Username: "test-user",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockLightClient()
			application := app.NewWithClient(cfg, mockClient)

			err := application.ExecuteCommand(tt.lightID, tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
