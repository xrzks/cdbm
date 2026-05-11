package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunCdCommand_StoreValidation(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	c := &CLI{store: s}
	
	// Test: Validate bookmark name directly (this is what cd command does first)
	err = c.store.ValidateBookmarkName("valid-name")
	if err != nil {
		t.Errorf("Unexpected error validating valid bookmark name: %v", err)
	}
	
	// Test: Validate invalid bookmark name
	err = c.store.ValidateBookmarkName("invalid name")
	if err == nil {
		t.Error("Expected error validating invalid bookmark name, but got none")
	}
	
	// Add a test bookmark
	testDir := filepath.Join(tmpDir, "cd-test")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	err = c.store.Add("cd-bookmark", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}
	
	// Test: Get bookmark (this is what cd command does after validation)
	bookmark, err := c.store.GetOne("cd-bookmark")
	if err != nil {
		t.Errorf("Failed to get bookmark: %v", err)
	}
	
	if bookmark.Name != "cd-bookmark" {
		t.Errorf("Expected bookmark name 'cd-bookmark', got '%s'", bookmark.Name)
	}
	
	if bookmark.Directory != testDir {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir, bookmark.Directory)
	}
	
	// Test: Get nonexistent bookmark
	_, err = c.store.GetOne("nonexistent")
	if err == nil {
		t.Error("Expected error getting nonexistent bookmark, but got none")
	}
}

func TestCLI_RunCdCommand_PathValidation(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	c := &CLI{store: s}
	
	// Test: Valid directory
	testDir := filepath.Join(tmpDir, "valid-dir")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create valid directory: %v", err)
	}
	
	err = c.store.Add("valid-dir", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark for valid directory: %v", err)
	}
	
	// Test: Invalid directory (nonexistent)
	err = c.store.Add("invalid-dir", "/nonexistent/path")
	if err == nil {
		t.Error("Expected error adding bookmark for nonexistent directory, but got none")
	}
	
	// Test: Invalid directory (file instead of directory)
	testFile := filepath.Join(tmpDir, "not-a-dir")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()
	
	err = c.store.Add("file-dir", testFile)
	if err == nil {
		t.Error("Expected error adding bookmark for file instead of directory, but got none")
	}
	
	// Test: Invalid directory (symlink)
	symlinkDir := filepath.Join(tmpDir, "symlink-target")
	err = os.MkdirAll(symlinkDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create symlink target: %v", err)
	}
	
	symlinkPath := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(symlinkDir, symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}
	
	err = c.store.Add("symlink-dir", symlinkPath)
	if err == nil {
		t.Error("Expected error adding bookmark for symlink directory, but got none")
	}
}

func TestCLI_RunCdCommand_PathCleaning(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	c := &CLI{store: s}
	
	// Test: Path with trailing slash (should be cleaned)
	testDir := filepath.Join(tmpDir, "clean-dir")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	// Add bookmark with trailing slash
	trailingSlashDir := testDir + "/"
	err = c.store.Add("clean-dir", trailingSlashDir)
	if err != nil {
		t.Errorf("Failed to add bookmark with trailing slash: %v", err)
	}
	
	// Get bookmark and verify path was cleaned
	bookmark, err := c.store.GetOne("clean-dir")
	if err != nil {
		t.Errorf("Failed to get bookmark: %v", err)
	}
	
	// The path should not have a trailing slash
	if bookmark.Directory != testDir {
		t.Errorf("Expected cleaned directory path '%s', got '%s'", testDir, bookmark.Directory)
	}
	
	// Test: Path with redundant components
	redundantPath := filepath.Join(tmpDir, "parent", "..", "child", ".", "final")
	cleanPath := filepath.Join(tmpDir, "child", "final")
	
	err = os.MkdirAll(cleanPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create redundant path directory: %v", err)
	}
	
	err = c.store.Add("redundant-dir", redundantPath)
	if err != nil {
		t.Errorf("Failed to add bookmark with redundant path: %v", err)
	}
	
	// Get bookmark and verify path was cleaned
	bookmark, err = c.store.GetOne("redundant-dir")
	if err != nil {
		t.Errorf("Failed to get bookmark: %v", err)
	}
	
	// The path should be cleaned
	if bookmark.Directory != cleanPath {
		t.Errorf("Expected cleaned directory path '%s', got '%s'", cleanPath, bookmark.Directory)
	}
}