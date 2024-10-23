package app

import (
	"fmt"
	"strings"

	"github.com/umit144/philips-hue-bulb-control/internal/config"
	"github.com/umit144/philips-hue-bulb-control/internal/light"
	"github.com/umit144/philips-hue-bulb-control/internal/ui"
)

type App struct {
	lightClient light.LightClient
	ui          *ui.UI
}

func New(cfg *config.Config) *App {
	lightClient := light.NewClient(cfg.BridgeIP, cfg.Username)
	return &App{
		lightClient: lightClient,
		ui:          ui.New(lightClient),
	}
}

func NewWithClient(cfg *config.Config, client light.LightClient) *App {
	return &App{
		lightClient: client,
		ui:          ui.New(client),
	}
}

func (a *App) StartInteractive() error {
	return a.ui.Run()
}

func (a *App) ExecuteCommand(lightID, stateStr string) error {
	lights, err := a.lightClient.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get lights: %w", err)
	}

	if _, exists := lights[lightID]; !exists {
		return fmt.Errorf("light with ID %s not found", lightID)
	}

	state := strings.ToLower(stateStr) == "on"
	if err := a.lightClient.Toggle(lightID, state); err != nil {
		return fmt.Errorf("failed to toggle light %s: %w", lightID, err)
	}

	fmt.Printf("Light %s turned %s\n", lightID, strings.ToUpper(stateStr))
	return nil
}
