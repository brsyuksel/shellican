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
	CreateCollection("col1")
	err = ListCollections()
	if err != nil {
		t.Errorf("ListCollections failed for one: %v", err)
	}
}

func TestListRunnables(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)

	CreateCollection("col1")
	// Add runnable to collection.yml
	colPath := filepath.Join(tempDir, ".shellican", "col1", "collection.yml")
	content := `
runnables:
  - run1
`
	os.WriteFile(colPath, []byte(content), 0644)

	// Create run dir
	runPath := filepath.Join(tempDir, ".shellican", "col1", "run1")
	os.MkdirAll(runPath, 0755)
	os.WriteFile(filepath.Join(runPath, "runnable.yml"), []byte("summary: test"), 0644)

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
