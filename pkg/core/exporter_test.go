package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportCollection(t *testing.T) {
	tempDir := t.TempDir()
	envHome := filepath.Join(tempDir, "env")
	t.Setenv("SHELLICAN_HOME", envHome)

	// Create collection to export
	colName := "export-test"
	colDir := filepath.Join(envHome, ".shellican", colName)
	if err := os.MkdirAll(colDir, 0755); err != nil {
		t.Fatalf("failed to create collection dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(colDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write content file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.tar.gz")

	// Tar might not act identical on all systems but typically available.
	err := ExportCollection(colName, outputFile)
	if err != nil {
		t.Fatalf("ExportCollection failed: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file not created")
	}

	// Test missing collection
	err = ExportCollection("missing", "out")
	if err == nil {
		t.Error("Expected error for missing collection")
	}
}
