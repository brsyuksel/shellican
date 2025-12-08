package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExportCollection exports a collection to a tar.gz file.
// name: name of the collection to export.
// output: optional output path. Defaults to <name>.tar.gz in current dir.
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

	// Ensure output path is absolute or relative to CWD.
	// tar command handles relative paths fine.

	fmt.Printf("Exporting collection '%s' to '%s'...\n", name, output)

	// tar -czf <output> -C <collectionPath> .
	// We export the CONTENTS of the collection directory so that when imported,
	// it can be extracted into a new directory cleanly?
	// Looking at import logic: `tar -xzf path -C target`.
	// If we export contents: `tar -czf out.tgz -C colDir .`
	// Then import: `mkdir target; tar -xzf out.tgz -C target` -> target/.
	// This matches our import logic.

	cmd := exec.Command("tar", "-czf", output, "-C", collectionPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar compression failed: %w", err)
	}

	fmt.Println("Export successful.")
	return nil
}
