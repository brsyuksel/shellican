package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brsyuksel/shellican/pkg/config"
)

func TestCreateCollection(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	colName := "test-col"
	err := CreateCollection(colName)
	if err != nil {
		t.Fatalf("CreateCollection failed: %v", err)
	}

	colPath := filepath.Join(tempDir, ".shellican", colName)
	if _, err := os.Stat(colPath); os.IsNotExist(err) {
		t.Errorf("Collection directory not created")
	}

	if _, err := os.Stat(filepath.Join(colPath, "collection.yml")); os.IsNotExist(err) {
		t.Errorf("collection.yml not created")
	}

	// Double create should fail
	err = CreateCollection(colName)
	if err == nil {
		t.Error("Expected error for existing collection, got nil")
	}
}

func TestCreateRunnable(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	colName := "test-col"
	if err := CreateCollection(colName); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	runName := "test-run"
	err := CreateRunnable(colName, runName)
	if err != nil {
		t.Fatalf("CreateRunnable failed: %v", err)
	}

	runPath := filepath.Join(tempDir, ".shellican", colName, runName)
	if _, err := os.Stat(runPath); os.IsNotExist(err) {
		t.Errorf("Runnable directory not created")
	}

	if _, err := os.Stat(filepath.Join(runPath, "runnable.yml")); os.IsNotExist(err) {
		t.Errorf("runnable.yml not created")
	}

	// Verify runnable was added to collection.yml
	colPath := filepath.Join(tempDir, ".shellican", colName)
	cfg, err := config.LoadCollectionConfig(colPath)
	if err != nil {
		t.Fatalf("Failed to load collection config: %v", err)
	}
	if cfg == nil {
		t.Fatal("Collection config is nil")
	}

	found := false
	for _, r := range cfg.Runnables {
		if r == runName {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Runnable '%s' not found in collection.yml runnables list", runName)
	}

	// Missing collection test
	err = CreateRunnable("missing", "run")
	if err == nil {
		t.Error("Expected error for missing collection, got nil")
	}
}
