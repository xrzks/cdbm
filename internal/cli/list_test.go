package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunListCommand(t *testing.T) {
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

	// Test: List with empty store
	bookmarks := c.store.GetAll()
	if len(bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks in empty store, got %d", len(bookmarks))
	}

	// Add some test bookmarks
	testDir1 := filepath.Join(tmpDir, "test1")
	testDir2 := filepath.Join(tmpDir, "test2")

	err = os.MkdirAll(testDir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory 1: %v", err)
	}

	err = os.MkdirAll(testDir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory 2: %v", err)
	}

	err = c.store.Add("bookmark1", testDir1)
	if err != nil {
		t.Errorf("Failed to add bookmark1: %v", err)
	}

	err = c.store.Add("bookmark2", testDir2)
	if err != nil {
		t.Errorf("Failed to add bookmark2: %v", err)
	}

	// Test: List with bookmarks
	bookmarks = c.store.GetAll()
	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}

	// Verify bookmark names
	foundNames := make(map[string]bool)
	for _, bookmark := range bookmarks {
		foundNames[bookmark.Name] = true
	}

	if !foundNames["bookmark1"] {
		t.Error("Expected bookmark1 not found")
	}

	if !foundNames["bookmark2"] {
		t.Error("Expected bookmark2 not found")
	}

	// Verify bookmark details
	bookmark1, err := c.store.GetOne("bookmark1")
	if err != nil {
		t.Errorf("Failed to get bookmark1: %v", err)
	}

	if bookmark1.Name != "bookmark1" {
		t.Errorf("Expected bookmark name 'bookmark1', got '%s'", bookmark1.Name)
	}

	if bookmark1.Directory != testDir1 {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir1, bookmark1.Directory)
	}
}

func TestCLI_RunListCommand_PrettyFormat(t *testing.T) {
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

	// Add a bookmark with spaces in the directory path
	testDir := filepath.Join(tmpDir, "test with spaces")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = c.store.Add("test-spaces", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}

	// Get the bookmark
	bookmark, err := c.store.GetOne("test-spaces")
	if err != nil {
		t.Errorf("Failed to get bookmark: %v", err)
	}

	// Test the Pretty method
	prettyOutput := bookmark.Pretty()

	// Verify the pretty output contains expected elements
	if prettyOutput == "" {
		t.Error("Pretty output is empty")
	}

	if !contains(prettyOutput, "test-spaces") {
		t.Error("Pretty output doesn't contain bookmark name")
	}

	if !contains(prettyOutput, testDir) {
		t.Error("Pretty output doesn't contain bookmark directory")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
