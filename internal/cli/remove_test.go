package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunRemoveCommand(t *testing.T) {
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

	// Test: Remove from empty store (should fail)
	err = c.store.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error removing from empty store, but got none")
	}

	// Add a test bookmark
	testDir := filepath.Join(tmpDir, "test-remove")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = c.store.Add("remove-test", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}

	// Verify bookmark exists
	_, err = c.store.GetOne("remove-test")
	if err != nil {
		t.Errorf("Failed to get bookmark before removal: %v", err)
	}

	// Test: Remove existing bookmark
	err = c.store.Delete("remove-test")
	if err != nil {
		t.Errorf("Failed to remove bookmark: %v", err)
	}

	// Verify bookmark is gone
	_, err = c.store.GetOne("remove-test")
	if err == nil {
		t.Error("Bookmark still exists after removal")
	}

	// Test: Remove nonexistent bookmark (should fail)
	err = c.store.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error removing nonexistent bookmark, but got none")
	}
}

func TestCLI_RunRemoveCommand_MultipleBookmarks(t *testing.T) {
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

	// Add multiple test bookmarks
	testDir1 := filepath.Join(tmpDir, "test1")
	testDir2 := filepath.Join(tmpDir, "test2")
	testDir3 := filepath.Join(tmpDir, "test3")

	dirs := []string{testDir1, testDir2, testDir3}
	names := []string{"bm1", "bm2", "bm3"}

	for i, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %d: %v", i, err)
		}

		err = c.store.Add(names[i], dir)
		if err != nil {
			t.Errorf("Failed to add bookmark %s: %v", names[i], err)
		}
	}

	// Verify all bookmarks exist
	for _, name := range names {
		_, err := c.store.GetOne(name)
		if err != nil {
			t.Errorf("Failed to get bookmark %s: %v", name, err)
		}
	}

	// Remove middle bookmark
	err = c.store.Delete("bm2")
	if err != nil {
		t.Errorf("Failed to remove bookmark bm2: %v", err)
	}

	// Verify removed bookmark is gone
	_, err = c.store.GetOne("bm2")
	if err == nil {
		t.Error("Bookmark bm2 still exists after removal")
	}

	// Verify other bookmarks still exist
	_, err = c.store.GetOne("bm1")
	if err != nil {
		t.Errorf("Failed to get bookmark bm1: %v", err)
	}

	_, err = c.store.GetOne("bm3")
	if err != nil {
		t.Errorf("Failed to get bookmark bm3: %v", err)
	}

	// Verify total count
	bookmarks := c.store.GetAll()
	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks after removal, got %d", len(bookmarks))
	}
}
