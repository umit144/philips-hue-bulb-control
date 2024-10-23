package ui_test

import (
	"os"
	"testing"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/umit144/philips-hue-bulb-control/internal/light"
	"github.com/umit144/philips-hue-bulb-control/internal/ui"
)

type MockLightClient struct {
	lights        map[string]light.Light
	getAllCalled  bool
	toggleCalled  bool
	toggleLightID string
	toggleState   bool
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
	m.getAllCalled = true
	return m.lights, nil
}

func (m *MockLightClient) Toggle(lightID string, state bool) error {
	m.toggleCalled = true
	m.toggleLightID = lightID
	m.toggleState = state
	return nil
}

func TestNewUI(t *testing.T) {
	mockClient := NewMockLightClient()
	ui := ui.New(mockClient)
	if ui == nil {
		t.Fatal("Expected UI instance, got nil")
	}
}

func TestDisplayLights(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping UI test in CI environment")
	}

	mockClient := NewMockLightClient()
	testUI := ui.New(mockClient)
	ui.SetTestMode(true)

	go func() {
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyEsc)
	}()

	err := testUI.Run()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !mockClient.getAllCalled {
		t.Error("Expected GetAll to be called")
	}
}

func TestLightToggling(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping UI test in CI environment")
	}

	mockClient := NewMockLightClient()
	testUI := ui.New(mockClient)
	ui.SetTestMode(true)

	go func() {
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyArrowDown)
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyEnter)
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyEsc)
	}()

	err := testUI.Run()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !mockClient.toggleCalled {
		t.Error("Expected Toggle to be called")
	}
	if mockClient.toggleLightID != "2" {
		t.Errorf("Expected light ID 2, got %s", mockClient.toggleLightID)
	}
	if mockClient.toggleState != true {
		t.Error("Expected toggle state to be true")
	}
}

func TestNavigateAndQuit(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping UI test in CI environment")
	}

	mockClient := NewMockLightClient()
	testUI := ui.New(mockClient)
	ui.SetTestMode(true)

	go func() {
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyArrowDown)
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyArrowUp)
		time.Sleep(100 * time.Millisecond)
		testUI.EmulateKeyPress(keyboard.KeyEsc)
	}()

	err := testUI.Run()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
