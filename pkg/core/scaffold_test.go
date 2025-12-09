package core

import (
	"os"
	"path/filepath"
	"testing"
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

	// Missing collection test
	err = CreateRunnable("missing", "run")
	if err == nil {
		t.Error("Expected error for missing collection, got nil")
	}
}
