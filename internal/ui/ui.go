package ui

import (
	"fmt"
	"sort"

	"github.com/eiannone/keyboard"
	"github.com/umit144/philips-hue-bulb-control/internal/light"
)

var (
	ErrStopped = fmt.Errorf("ui stopped")
	testMode   = false
)

type UI struct {
	lightClient light.LightClient
	stopped     bool
	keyEvents   chan keyboard.Key
	charEvents  chan rune
}

func New(lightClient light.LightClient) *UI {
	return &UI{
		lightClient: lightClient,
		keyEvents:   make(chan keyboard.Key, 10),
		charEvents:  make(chan rune, 10),
	}
}

func SetTestMode(enabled bool) {
	testMode = enabled
}

func (ui *UI) Stop() {
	ui.stopped = true
}

func (ui *UI) EmulateKeyPress(key keyboard.Key) {
	ui.keyEvents <- key
}

func (ui *UI) EmulateChar(char rune) {
	ui.charEvents <- char
}

func (ui *UI) Run() error {
	if err := keyboard.Open(); err != nil {
		return fmt.Errorf("failed to open keyboard: %w", err)
	}
	defer keyboard.Close()

	selectedIndex := 0

	for {
		if ui.stopped {
			return ErrStopped
		}

		lights, err := ui.lightClient.GetAll()
		if err != nil {
			return fmt.Errorf("failed to get lights: %w", err)
		}

		sortedLights := ui.getSortedLights(lights)
		ui.display(sortedLights, selectedIndex)

		if err := ui.handleInput(&selectedIndex, sortedLights); err != nil {
			if err == ErrStopped {
				return nil
			}
			return fmt.Errorf("input handling failed: %w", err)
		}
	}
}

func (ui *UI) handleInput(selectedIndex *int, lights []light.Light) error {
	select {
	case key := <-ui.keyEvents:
		return ui.handleKey(key, selectedIndex, lights)
	case char := <-ui.charEvents:
		if char == KeyQuit {
			return ErrStopped
		}
	default:
		char, key, err := keyboard.GetKey()
		if err != nil {
			return fmt.Errorf("failed to read keyboard: %w", err)
		}
		if char != 0 {
			ui.charEvents <- char
		} else {
			ui.keyEvents <- key
		}
	}

	return nil
}

func (ui *UI) handleKey(key keyboard.Key, selectedIndex *int, lights []light.Light) error {
	switch key {
	case KeyArrowUp:
		if *selectedIndex > 0 {
			*selectedIndex--
		}
	case KeyArrowDown:
		if *selectedIndex < len(lights)-1 {
			*selectedIndex++
		}
	case KeyEnter:
		if *selectedIndex < len(lights) {
			selected := lights[*selectedIndex]
			if err := ui.lightClient.Toggle(selected.ID, !selected.State.On); err != nil {
				return fmt.Errorf("failed to toggle light: %w", err)
			}
		}
	case keyboard.KeyEsc:
		return ErrStopped
	}
	return nil
}

func (ui *UI) getSortedLights(lights map[string]light.Light) []light.Light {
	sorted := make([]light.Light, 0, len(lights))
	for _, l := range lights {
		sorted = append(sorted, l)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})

	return sorted
}

func (ui *UI) display(lights []light.Light, selectedIndex int) {
	if !testMode {
		fmt.Print("\033[H\033[2J")
		fmt.Println("\n=== Philips Hue Light Control ===")
		fmt.Println("Use ↑↓ arrows to navigate, Enter to toggle, 'q' to quit")
		fmt.Println("=====================================")
	}

	for i, light := range lights {
		cursor := "  "
		if i == selectedIndex {
			cursor = "→ "
		}

		state := "OFF"
		if light.State.On {
			state = "ON "
		}

		if !testMode {
			fmt.Printf("%s%s\t%-20s\t[%s]\n", cursor, light.ID, light.Name, state)
		}
	}
}
