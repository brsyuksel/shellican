package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// CreateShell creates a wrapper script for the collection.
func CreateShell(collection, name string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

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

	helperName := name
	if helperName == "" {
		helperName = fmt.Sprintf("%s-shell", collection)
	}

	if runtime.GOOS == "windows" {
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
