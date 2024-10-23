package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/umit144/philips-hue-bulb-control/internal/config"
)

func TestLoad(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	envContent := `
	HUE_BRIDGE_IP=192.168.1.100
	HUE_USERNAME=testuser
	`

	if err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to write test .env file: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.BridgeIP != "192.168.1.100" {
		t.Errorf("Expected BridgeIP=192.168.1.100, got %s", cfg.BridgeIP)
	}

	if cfg.Username != "testuser" {
		t.Errorf("Expected Username=testuser, got %s", cfg.Username)
	}
}
