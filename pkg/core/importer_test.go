package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestImportCollection_Folder(t *testing.T) {
	tempDir := t.TempDir()
	envHome := filepath.Join(tempDir, "env")
	t.Setenv("SHELLICAN_HOME", envHome)

	// Prepare source
	sourceDir := filepath.Join(tempDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("data"), 0644)

	err := ImportCollection(sourceDir, "imported-folder")
	if err != nil {
		t.Fatalf("ImportCollection folder failed: %v", err)
	}

	resultDir := filepath.Join(envHome, ".shellican", "imported-folder")
	if _, err := os.Stat(filepath.Join(resultDir, "test.txt")); os.IsNotExist(err) {
		t.Errorf("Imported file missing")
	}
}

func TestImportCollection_Tarball(t *testing.T) {
	tempDir := t.TempDir()
	envHome := filepath.Join(tempDir, "env")
	t.Setenv("SHELLICAN_HOME", envHome)

	// Prepare tarball
	contentDir := filepath.Join(tempDir, "content")
	os.MkdirAll(contentDir, 0755)
	os.WriteFile(filepath.Join(contentDir, "data.txt"), []byte("tar-data"), 0644)

	tarFile := filepath.Join(tempDir, "test.tar.gz")
	cmd := exec.Command("tar", "-czf", tarFile, "-C", contentDir, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create test tarball: %v", err)
	}

	err := ImportCollection(tarFile, "imported-tar")
	if err != nil {
		t.Fatalf("ImportCollection tarball failed: %v", err)
	}

	resultDir := filepath.Join(envHome, ".shellican", "imported-tar")
	if _, err := os.Stat(filepath.Join(resultDir, "data.txt")); os.IsNotExist(err) {
		t.Errorf("Imported tar file missing")
	}
}
