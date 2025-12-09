package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CreateShellHelper creates a wrapper script for the collection.
// Ideally this should go to a location in PATH.
// For now, we will try to put it in $HOME/bin or similar if it exists,
// or just tell the user where we created it.
func CreateShell(collection, name string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

	// Determine target binary path
	// Priority: $HOME/bin, $HOME/.local/bin, or just in current dir?
	// User request: "bin path'e bir helper shell atilacak"
	// Let's try standard user bin paths.

	binPaths := []string{
		filepath.Join(homeDir, ".local", "bin"),
	}

	var targetDir string
	for _, p := range binPaths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			targetDir = p
			break
		}
	}

	if targetDir == "" {
		return fmt.Errorf("could not find a suitable bin directory (~/.local/bin). Please create it and add to PATH")
	}

	// The helper should call `shellican <collection> "$@"`
	// We assume `shellican` is in the PATH as well.

	helperName := name
	if helperName == "" {
		helperName = fmt.Sprintf("%s-shell", collection)
	}

	if runtime.GOOS == "windows" {
		// Windows support removed for this version
		return fmt.Errorf("windows support is disabled in this version")
	}

	targetPath := filepath.Join(targetDir, helperName)

	content := fmt.Sprintf("#!/bin/sh\nexec shellican run %s \"$@\"\n", collection)

	err = os.WriteFile(targetPath, []byte(content), 0755)
	if err != nil {
		return fmt.Errorf("failed to write helper script to %s: %w", targetPath, err)
	}

	fmt.Printf("Helper created at: %s\n", targetPath)
	return nil
}
