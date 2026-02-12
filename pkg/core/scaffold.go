package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brsyuksel/shellican/pkg/config"
)

// CreateCollection creates a new collection.
func CreateCollection(name string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	collectionPath := filepath.Join(rootDir, name)
	if _, err := os.Stat(collectionPath); !os.IsNotExist(err) {
		return fmt.Errorf("collection already exists: %s", collectionPath)
	}

	if err := os.MkdirAll(collectionPath, 0755); err != nil {
		return fmt.Errorf("failed to create collection directory: %w", err)
	}

	configContent := fmt.Sprintf(`name: "%s"
help: "Usage for %s"
readme: "README.md"
runnables: []
environments:
  COLLECTION_ENV: "value"
`, name, name)

	configPath := filepath.Join(collectionPath, "collection.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write collection.yml: %w", err)
	}

	readmeContent := fmt.Sprintf("# %s\n\nDescription for %s\n", name, name)
	readmePath := filepath.Join(collectionPath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		fmt.Printf("Warning: failed to create README.md: %v\n", err)
	}

	fmt.Printf("Collection '%s' created at %s\n", name, collectionPath)
	return nil
}

// CreateRunnable creates a new runnable.
func CreateRunnable(collectionName, runnableName string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	collectionPath := filepath.Join(rootDir, collectionName)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return fmt.Errorf("collection does not exist: %s, please create it first", collectionPath)
	}

	runnablePath := filepath.Join(collectionPath, runnableName)
	if _, err := os.Stat(runnablePath); !os.IsNotExist(err) {
		return fmt.Errorf("runnable already exists: %s", runnablePath)
	}

	if err := os.MkdirAll(runnablePath, 0755); err != nil {
		return fmt.Errorf("failed to create runnable directory: %w", err)
	}

	configContent := fmt.Sprintf(`name: "%s"
help: "Usage for %s"
readme: "README.md"
run: "echo 'Hello from %s'"
# before: "echo 'Running before'"
# after: "echo 'Running after'"
environments:
  RUNNABLE_ENV: "value"
`, runnableName, runnableName, runnableName)

	configPath := filepath.Join(runnablePath, "runnable.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write runnable.yml: %w", err)
	}

	readmeContent := fmt.Sprintf("# %s\n\nDescription for %s\n", runnableName, runnableName)
	readmePath := filepath.Join(runnablePath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		fmt.Printf("Warning: failed to create README.md: %v\n", err)
	}

	// Load collection config and add the runnable
	collectionCfg, err := config.LoadCollectionConfig(collectionPath)
	if err != nil {
		return fmt.Errorf("failed to load collection config: %w", err)
	}
	if collectionCfg == nil {
		return fmt.Errorf("collection config not found")
	}

	// Add runnable to the list if not already present
	found := false
	for _, r := range collectionCfg.Runnables {
		if r == runnableName {
			found = true
			break
		}
	}
	if !found {
		collectionCfg.Runnables = append(collectionCfg.Runnables, runnableName)
	}

	// Save updated collection config
	if err := config.SaveCollectionConfig(collectionPath, collectionCfg); err != nil {
		return fmt.Errorf("failed to save collection config: %w", err)
	}

	fmt.Printf("Runnable '%s' created and added to collection at %s\n", runnableName, runnablePath)
	return nil
}
