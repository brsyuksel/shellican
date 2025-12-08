package core

import (
	"testing"
)

func TestShowCollection(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	CreateCollection("col1")
	err := ShowCollection("col1", false)
	if err != nil {
		t.Errorf("ShowCollection failed: %v", err)
	}

	err = ShowCollection("missing", false)
	if err == nil {
		t.Error("Expected error for missing collection")
	}
}

func TestShowRunnable(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	CreateCollection("col1")
	// Manually construct proper runnable setup since CreateRunnable does minimal config
	CreateRunnable("col1", "run1")

	err := ShowRunnable("col1", "run1", false)
	if err != nil {
		t.Errorf("ShowRunnable failed: %v", err)
	}

	err = ShowRunnable("col1", "missing", false)
	if err == nil {
		t.Error("Expected error for missing runnable")
	}
}
