package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateShell(t *testing.T) {
	tempDir := t.TempDir()
	// Mock HOME to point to tempDir so CreateShell finds/creates .local/bin there
	t.Setenv("HOME", tempDir)

	// Create .local/bin
	binDir := filepath.Join(tempDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin dir: %v", err)
	}

	colName := "test-col"
	err := CreateShell(colName, "")
	if err != nil {
		t.Fatalf("CreateShell failed: %v", err)
	}

	expectedScript := filepath.Join(binDir, colName+"-shell")
	if _, err := os.Stat(expectedScript); os.IsNotExist(err) {
		t.Errorf("Expected shell script at %s, not found", expectedScript)
	}

	// Test with custom name
	customName := "my-helper"
	err = CreateShell(colName, customName)
	if err != nil {
		t.Fatalf("CreateShell with name failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(binDir, customName)); os.IsNotExist(err) {
		t.Errorf("Expected custom shell script at %s, not found", customName)
	}
}
