package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ImportCollection imports a collection from a source.
func ImportCollection(source, name string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	if name == "" {
		base := filepath.Base(source)
		if strings.HasSuffix(base, ".git") {
			name = strings.TrimSuffix(base, ".git")
		} else if strings.HasSuffix(base, ".tar.gz") {
			name = strings.TrimSuffix(base, ".tar.gz")
		} else if strings.HasSuffix(base, ".tgz") {
			name = strings.TrimSuffix(base, ".tgz")
		} else {
			name = base
		}
	}

	targetDir := filepath.Join(rootDir, name)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return fmt.Errorf("collection '%s' already exists at %s", name, targetDir)
	}

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "git@") || strings.HasSuffix(source, ".git") {
		return importGit(source, targetDir)
	}

	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return importTarball(source, targetDir)
	}

	info, err := os.Stat(source)
	if err == nil && info.IsDir() {
		return importFolder(source, targetDir)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error checking source: %w", err)
	}

	return fmt.Errorf("unknown source type or source not found: %s", source)
}

// importGit clones a git repository.
func importGit(url, target string) error {
	fmt.Printf("Importing from Git: %s...\n", url)
	cmd := exec.Command("git", "clone", url, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}

// importTarball extracts a tarball.
func importTarball(path, target string) error {
	fmt.Printf("Importing from Tarball: %s...\n", path)
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create target dir: %w", err)
	}

	cmd := exec.Command("tar", "-xzf", path, "-C", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(target) // Cleanup
		return fmt.Errorf("tar extraction failed: %w", err)
	}
	return nil
}

// importFolder copies a folder.
func importFolder(source, target string) error {
	fmt.Printf("Importing from Folder: %s...\n", source)

	fmt.Printf("Importing from Folder: %s...\n", source)

	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("failed to clean target path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("failed to create parent dir: %w", err)
	}

	cmd := exec.Command("cp", "-R", source, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("folder copy failed: %w", err)
	}
	return nil
}
