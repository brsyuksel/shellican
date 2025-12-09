package core

import (
	"testing"
)

func TestShowCollection(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	if err := CreateCollection("col1"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
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

	if err := CreateCollection("col1"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	// Manually construct proper runnable setup since CreateRunnable does minimal config
	if err := CreateRunnable("col1", "run1"); err != nil {
		t.Fatalf("setup run failed: %v", err)
	}

	err := ShowRunnable("col1", "run1", false)
	if err != nil {
		t.Errorf("ShowRunnable failed: %v", err)
	}

	err = ShowRunnable("col1", "missing", false)
	if err == nil {
		t.Error("Expected error for missing runnable")
	}
}
