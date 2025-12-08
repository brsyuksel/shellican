package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCollectionConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "collection.yml")
	content := `
summary: "Test Collection"
help: "This is a test"
runnables:
  - test-runnable
environments:
  KEY: VALUE
`
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loading from directory (searching for collection.yml)
	cfg, err := LoadCollectionConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load collection config: %v", err)
	}
	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}

	if cfg.Summary != "Test Collection" {
		t.Errorf("Expected summary 'Test Collection', got '%s'", cfg.Summary)
	}
	if len(cfg.Runnables) != 1 || cfg.Runnables[0] != "test-runnable" {
		t.Errorf("Runnables mismatch")
	}
	if cfg.Environments["KEY"] != "VALUE" {
		t.Errorf("Environment variable mismatch")
	}
}

func TestLoadRunnableConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "runnable.yml")
	content := `
summary: "Test Runnable"
type: "inline"
run: "echo hello"
environments:
  RUN_KEY: RUN_VALUE
`
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadRunnableConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load runnable config: %v", err)
	}

	if cfg.Summary != "Test Runnable" {
		t.Errorf("Expected summary 'Test Runnable', got '%s'", cfg.Summary)
	}
	if cfg.Type != "inline" {
		t.Errorf("Expected type 'inline', got '%s'", cfg.Type)
	}
	if cfg.Environments["RUN_KEY"] != "RUN_VALUE" {
		t.Errorf("Environment variable mismatch")
	}
}
