package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExportCollection exports a collection to a tar.gz file.
func ExportCollection(name, output string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	collectionPath := filepath.Join(rootDir, name)
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return fmt.Errorf("collection '%s' does not exist", name)
	}

	if output == "" {
		output = fmt.Sprintf("%s.tar.gz", name)
	}

	fmt.Printf("Exporting collection '%s' to '%s'...\n", name, output)

	cmd := exec.Command("tar", "-czf", output, "-C", collectionPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar compression failed: %w", err)
	}

	fmt.Println("Export successful.")
	return nil
}
