package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func getRoot() (string, error) {
	if envHome := os.Getenv("SHELLICAN_HOME"); envHome != "" {
		return filepath.Join(envHome, ".shellican"), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}
	return filepath.Join(homeDir, ".shellican"), nil
}
