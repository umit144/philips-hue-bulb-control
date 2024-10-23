package mocks

import (
	"fmt"

	"github.com/umit144/philips-hue-bulb-control/internal/light"
)

type MockLightClient struct {
	lights map[string]light.Light
	err    error
}

func NewMockLightClient() *MockLightClient {
	return &MockLightClient{
		lights: map[string]light.Light{
			"1": {ID: "1", Name: "Living Room", State: light.State{On: true}},
			"2": {ID: "2", Name: "Kitchen", State: light.State{On: false}},
		},
	}
}

func (m *MockLightClient) GetAll() (map[string]light.Light, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.lights, nil
}

func (m *MockLightClient) Toggle(lightID string, state bool) error {
	if m.err != nil {
		return m.err
	}
	if light, exists := m.lights[lightID]; exists {
		light.State.On = state
		m.lights[lightID] = light
		return nil
	}
	return fmt.Errorf("light with ID %s not found", lightID)
}

func (m *MockLightClient) SetError(err error) {
	m.err = err
}
