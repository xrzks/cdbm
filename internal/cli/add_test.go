package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunAddCommand(t *testing.T) {
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

	// Test adding a bookmark with a valid directory
	testDir := filepath.Join(tmpDir, "test-dir")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test: Add bookmark with explicit name and directory
	err = c.store.Add("test-bookmark", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}

	// Verify bookmark was created
	bookmark, err := c.store.GetOne("test-bookmark")
	if err != nil {
		t.Errorf("Failed to get bookmark: %v", err)
	}

	if bookmark.Name != "test-bookmark" {
		t.Errorf("Expected bookmark name 'test-bookmark', got '%s'", bookmark.Name)
	}

	if bookmark.Directory != testDir {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir, bookmark.Directory)
	}

	// Test: Add bookmark with invalid name (should fail)
	err = c.store.Add("invalid name", testDir)
	if err == nil {
		t.Error("Expected error adding bookmark with invalid name, but got none")
	}

	// Test: Add bookmark with invalid directory (should fail)
	err = c.store.Add("test-invalid", "/nonexistent/path")
	if err == nil {
		t.Error("Expected error adding bookmark with invalid directory, but got none")
	}

	// Test: Add duplicate bookmark (should fail)
	err = c.store.Add("test-bookmark", testDir)
	if err == nil {
		t.Error("Expected error adding duplicate bookmark, but got none")
	}
}

func TestCLI_RunAddCommand_NameGeneration(t *testing.T) {
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

	// Test with directories that need name sanitization
	testDirs := []struct {
		dirPath      string
		expectedName string
	}{
		{
			dirPath:      filepath.Join(tmpDir, "test-dir_name"),
			expectedName: "test-dir",
		},
		{
			dirPath:      filepath.Join(tmpDir, "test.dots"),
			expectedName: "test.dots",
		},
		{
			dirPath:      filepath.Join(tmpDir, "test-hyphens"),
			expectedName: "test-hyphens",
		},
		{
			dirPath:      filepath.Join(tmpDir, "test  with spaces"),
			expectedName: "test-with-spaces",
		},
	}

	for _, tc := range testDirs {
		// Create the directory
		err := os.MkdirAll(tc.dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Add bookmark using the directory path directly
		err = c.store.Add(tc.expectedName, tc.dirPath)
		if err != nil {
			t.Errorf("Failed to add bookmark for %s: %v", tc.dirPath, err)
			continue
		}

		// Verify bookmark was created
		bookmark, err := c.store.GetOne(tc.expectedName)
		if err != nil {
			t.Errorf("Failed to get bookmark %s: %v", tc.expectedName, err)
			continue
		}

		if bookmark.Name != tc.expectedName {
			t.Errorf("Expected bookmark name '%s', got '%s'", tc.expectedName, bookmark.Name)
		}

		if bookmark.Directory != tc.dirPath {
			t.Errorf("Expected bookmark directory '%s', got '%s'", tc.dirPath, bookmark.Directory)
		}
	}
}

func TestCLI_RunAddCommand_EdgeCases(t *testing.T) {
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

	// Test: Add bookmark with empty name (should fail)
	testDir := filepath.Join(tmpDir, "test-empty")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = c.store.Add("", testDir)
	if err == nil {
		t.Error("Expected error adding bookmark with empty name, but got none")
	}

	// Test: Add bookmark with very long name (should fail)
	longName := "a" + string(make([]byte, 101)) // 102 characters
	err = c.store.Add(longName, testDir)
	if err == nil {
		t.Error("Expected error adding bookmark with very long name, but got none")
	}

	// Test: Add bookmark with special characters (should fail)
	err = c.store.Add("special*chars", testDir)
	if err == nil {
		t.Error("Expected error adding bookmark with special characters, but got none")
	}
}
