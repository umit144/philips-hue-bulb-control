package light

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type State struct {
	On  bool `json:"on"`
	Bri int  `json:"bri,omitempty"`
	Hue int  `json:"hue,omitempty"`
	Sat int  `json:"sat,omitempty"`
}

type Light struct {
	ID    string `json:"-"`
	State State  `json:"state"`
	Name  string `json:"name"`
	Type  string `json:"type,omitempty"`
}

type Client struct {
	bridgeIP string
	username string
	client   *http.Client
	baseURL  string
}

func NewClient(bridgeIP, username string) *Client {
	return &Client{
		bridgeIP: bridgeIP,
		username: username,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) getURL(path string) string {
	if c.baseURL != "" {
		return c.baseURL + path
	}
	return fmt.Sprintf("https://%s/api/%s%s", c.bridgeIP, c.username, path)
}

func (c *Client) GetAll() (map[string]Light, error) {
	url := c.getURL("/lights")

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get lights: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawLights map[string]Light
	if err := json.NewDecoder(resp.Body).Decode(&rawLights); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := make(map[string]Light)
	for id, light := range rawLights {
		light.ID = id
		result[id] = light
	}

	return result, nil
}

func (c *Client) Toggle(lightID string, state bool) error {
	url := c.getURL(fmt.Sprintf("/lights/%s/state", lightID))

	data, err := json.Marshal(State{On: state})
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, item := range result {
		if _, ok := item["success"]; ok {
			return nil
		}
		if errInfo, ok := item["error"]; ok {
			return fmt.Errorf("API error: %v", errInfo)
		}
	}

	return nil
}
