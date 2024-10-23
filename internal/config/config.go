package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BridgeIP string
	Username string
}

func Load() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	envPath := filepath.Join(currentDir, ".env")

	file, err := os.Open(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		switch key {
		case "HUE_BRIDGE_IP":
			config.BridgeIP = value
		case "HUE_USERNAME":
			config.Username = value
		}
	}

	if config.BridgeIP == "" || config.Username == "" {
		return nil, fmt.Errorf("missing required configuration")
	}

	return config, nil
}
