package main

import (
	"log"
	"os"

	"github.com/umit144/philips-hue-bulb-control/internal/app"
	"github.com/umit144/philips-hue-bulb-control/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application := app.New(cfg)

	if len(os.Args) == 3 {
		if err := application.ExecuteCommand(os.Args[1], os.Args[2]); err != nil {
			log.Fatalf("Command execution failed: %v", err)
		}
		return
	}

	if err := application.StartInteractive(); err != nil {
		log.Fatalf("Interactive mode failed: %v", err)
	}
}
