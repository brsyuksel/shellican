package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brsyuksel/shellican/pkg/config"
)

// ShowCollection prints details about a collection
func ShowCollection(name string, showReadme bool) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}
	collectionPath := filepath.Join(rootDir, name)

	cfg, err := config.LoadCollectionConfig(collectionPath)
	if err != nil {
		return fmt.Errorf("failed to load collection config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("collection '%s' not found", name)
	}

	fmt.Printf("Collection: %s\n", name)
	fmt.Printf("Summary:    %s\n", cfg.Summary)
	fmt.Printf("Help:       %s\n", cfg.Help)

	if showReadme && cfg.Readme != "" {
		readmePath := filepath.Join(collectionPath, cfg.Readme)
		content, err := os.ReadFile(readmePath)
		if err != nil {
			fmt.Printf("Warning: Failed to read README at %s: %v\n", readmePath, err)
		} else {
			fmt.Println("\n--- README ---")
			fmt.Println(string(content))
		}
	} else if showReadme {
		fmt.Println("No readme specified in configuration.")
	}

	return nil
}

// ShowRunnable prints details about a runnable
func ShowRunnable(collectionName, runnableName string, showReadme bool) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}
	collectionPath := filepath.Join(rootDir, collectionName)
	runnablePath := filepath.Join(collectionPath, runnableName)

	cfg, err := config.LoadRunnableConfig(runnablePath)
	if err != nil {
		return fmt.Errorf("failed to load runnable config: %w", err)
	}
	if cfg == nil {
		// Try to see if it is just unlisted or really non-existent
		if _, err := os.Stat(runnablePath); os.IsNotExist(err) {
			return fmt.Errorf("runnable '%s' not found in collection '%s'", runnableName, collectionName)
		}
		return fmt.Errorf("runnable found but invalid config")
	}

	fmt.Printf("Runnable:   %s (Collection: %s)\n", runnableName, collectionName)
	fmt.Printf("Summary:    %s\n", cfg.Summary)
	fmt.Printf("Help:       %s\n", cfg.Help)
	fmt.Printf("Type:       %s\n", cfg.Type)
	if cfg.Type == "inline" {
		fmt.Printf("Run Code:   %s\n", cfg.Run)
	} else {
		fmt.Printf("Run Script: %s\n", cfg.Run)
	}

	if showReadme && cfg.Readme != "" {
		readmePath := filepath.Join(runnablePath, cfg.Readme)
		content, err := os.ReadFile(readmePath)
		if err != nil {
			fmt.Printf("Warning: Failed to read README at %s: %v\n", readmePath, err)
		} else {
			fmt.Println("\n--- README ---")
			fmt.Println(string(content))
		}
	} else if showReadme {
		fmt.Println("No readme specified in configuration.")
	}

	return nil
}
