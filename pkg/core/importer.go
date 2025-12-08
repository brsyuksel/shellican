package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ImportCollection imports a collection from a source.
// source: git URL, tar.gz path, or local folder path.
// name: optional name for the collection. If empty, inferred from source.
func ImportCollection(source, name string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	// Infer name if not provided
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

	// 1. Check for Git
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "git@") || strings.HasSuffix(source, ".git") {
		return importGit(source, targetDir)
	}

	// 2. Check for Tarball
	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return importTarball(source, targetDir)
	}

	// 3. Check for Local Folder
	info, err := os.Stat(source)
	if err == nil && info.IsDir() {
		return importFolder(source, targetDir)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error checking source: %w", err)
	}

	return fmt.Errorf("unknown source type or source not found: %s", source)
}

func importGit(url, target string) error {
	fmt.Printf("Importing from Git: %s...\n", url)
	cmd := exec.Command("git", "clone", url, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	// Remove .git directory to make it a pure collection?
	// Usually user might want updates. Let's keep .git for now unless specified otherwise.
	// Spec didn't say. Keeping it is safer for "import from git" usually.
	return nil
}

func importTarball(path, target string) error {
	fmt.Printf("Importing from Tarball: %s...\n", path)
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create target dir: %w", err)
	}

	// Extract to target
	// tar -xzf <path> -C <target>
	cmd := exec.Command("tar", "-xzf", path, "-C", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(target) // Cleanup
		return fmt.Errorf("tar extraction failed: %w", err)
	}
	return nil
}

func importFolder(source, target string) error {
	fmt.Printf("Importing from Folder: %s...\n", source)

	// Ensure target doesn't exist, as we want to copy source AS target
	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("failed to clean target path: %w", err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("failed to create parent dir: %w", err)
	}

	// cp -R source target
	// Since target does not exist, this copies source directory TO target path.
	cmd := exec.Command("cp", "-R", source, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("folder copy failed: %w", err)
	}
	return nil
}
