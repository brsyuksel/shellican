package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brsyuksel/shellican/pkg/config"
)

// Mock helpers or setup environment for testing ResolveCommand would be complex
// because it relies on getRoot() which reads env vars or home dir.
// For this unit test, we will try to set SHELLICAN_HOME to a temp directory.

func TestResolveCommand(t *testing.T) {
	// Setup temporary SHELLICAN_HOME
	tempDir := t.TempDir()
	t.Setenv("SHELLICAN_HOME", tempDir)
	shellicanDir := filepath.Join(tempDir, ".shellican")
	if err := os.MkdirAll(shellicanDir, 0755); err != nil {
		t.Fatalf("Failed to create .shellican dir: %v", err)
	}

	// Create a test collection
	colName := "test-col"
	colDir := filepath.Join(shellicanDir, colName)
	if err := os.MkdirAll(colDir, 0755); err != nil {
		t.Fatalf("Failed to create collection dir: %v", err)
	}

	// Create collection.yml
	colContent := `
runnables:
  - my-script
environments:
  COL_ENV: "col"
`
	if err := os.WriteFile(filepath.Join(colDir, "collection.yml"), []byte(colContent), 0644); err != nil {
		t.Fatalf("Failed to write collection.yml: %v", err)
	}

	// Create runnable
	runName := "my-script"
	runDir := filepath.Join(colDir, runName)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		t.Fatalf("Failed to create runnable dir: %v", err)
	}

	// Create runnable.yml
	runContent := `
type: inline
run: echo hi
environments:
  RUN_ENV: "run"
`
	if err := os.WriteFile(filepath.Join(runDir, "runnable.yml"), []byte(runContent), 0644); err != nil {
		t.Fatalf("Failed to write runnable.yml: %v", err)
	}

	// Test valid resolution
	ctx, err := ResolveCommand(colName, []string{runName})
	if err != nil {
		t.Fatalf("ResolveCommand failed: %v", err)
	}

	if ctx.Config.Type != "inline" {
		t.Errorf("Expected inline type, got %s", ctx.Config.Type)
	}
	if ctx.Environments["COL_ENV"] != "col" {
		t.Errorf("Missing collection env")
	}
	if ctx.Environments["RUN_ENV"] != "run" {
		t.Errorf("Missing runnable env")
	}

	// Test invalid resolution (nested path invalid now)
	_, err = ResolveCommand(colName, []string{runName, "nested"})
	if err == nil {
		t.Error("Expected error for nested path components, got nil")
	}

	// Test unlisted runnable
	hiddenRunName := "hidden"
	hiddenRunDir := filepath.Join(colDir, hiddenRunName)
	os.MkdirAll(hiddenRunDir, 0755)
	os.WriteFile(filepath.Join(hiddenRunDir, "runnable.yml"), []byte("type: inline"), 0644)

	_, err = ResolveCommand(colName, []string{hiddenRunName})
	if err == nil {
		t.Error("Expected error for unlisted runnable, got nil")
	}
}

func TestExecuteContext(t *testing.T) {
	// Simple execution test
	// We can't easily capture stdout without redirecting os.Stdout,
	// but we can check if it returns valid error or nil.

	ctx := &ExecutionContext{
		RunnablePath: "/tmp", // dummy
		Config: &config.RunnableConfig{
			Type: "inline",
			Run:  "true",
		},
		Environments: map[string]string{},
	}

	if err := ExecuteContext(ctx, []string{}); err != nil {
		t.Errorf("ExecuteContext failed for simple true command: %v", err)
	}
}
