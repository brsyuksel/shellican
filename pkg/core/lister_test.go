package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListCollections(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	// No collections
	err := ListCollections()
	if err != nil {
		t.Errorf("ListCollections failed for empty: %v", err)
	}

	// One collection
	if err := CreateCollection("col1"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	err = ListCollections()
	if err != nil {
		t.Errorf("ListCollections failed for one: %v", err)
	}
}

func TestListRunnables(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	if err := CreateCollection("col1"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	// Add runnable to collection.yml
	colPath := filepath.Join(tempDir, ".shellican", "col1", "collection.yml")
	content := `
runnables:
  - run1
`
	if err := os.WriteFile(colPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write collection config: %v", err)
	}

	// Create run dir
	runPath := filepath.Join(tempDir, ".shellican", "col1", "run1")
	if err := os.MkdirAll(runPath, 0755); err != nil {
		t.Fatalf("failed to create run dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runPath, "runnable.yml"), []byte("summary: test"), 0644); err != nil {
		t.Fatalf("failed to write runnable config: %v", err)
	}

	err := ListRunnables("col1")
	if err != nil {
		t.Errorf("ListRunnables failed: %v", err)
	}

	// Bad collection
	err = ListRunnables("missing")
	if err == nil {
		t.Error("Expected error for missing collection")
	}
}
