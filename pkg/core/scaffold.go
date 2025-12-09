package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateCollection creates a new collection directory and a default collection.yml
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

	configContent := fmt.Sprintf(`name: "Description for %s"
help: "Usage for %s"
readme: "README.md"
runnables: []
  # - example_runnable
`, name, name)

	configPath := filepath.Join(collectionPath, "collection.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write collection.yml: %w", err)
	}

	fmt.Printf("Collection '%s' created at %s\n", name, collectionPath)
	return nil
}

// CreateRunnable creates a new runnable directory and a default runnable.yml within a collection
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

	configContent := fmt.Sprintf(`name: "Description for %s"
help: "Usage for %s"
type: "inline"
run: "echo 'Hello from %s'"
`, runnableName, runnableName, runnableName)

	configPath := filepath.Join(runnablePath, "runnable.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write runnable.yml: %w", err)
	}

	fmt.Printf("Runnable '%s' created at %s\n", runnableName, runnablePath)
	fmt.Println("Hint: Don't forget to add it to 'runnables' in collection.yml!")
	return nil
}
