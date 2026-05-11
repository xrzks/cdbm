package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunEditCommand(t *testing.T) {
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
	
	// Add a test bookmark
	testDir1 := filepath.Join(tmpDir, "original")
	testDir2 := filepath.Join(tmpDir, "updated")
	
	err = os.MkdirAll(testDir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory 1: %v", err)
	}
	
	err = os.MkdirAll(testDir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory 2: %v", err)
	}
	
	err = c.store.Add("original-name", testDir1)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}
	
	// Test: Edit name only
	err = c.store.Edit("original-name", "edited-name", "")
	if err != nil {
		t.Errorf("Failed to edit bookmark name: %v", err)
	}
	
	// Verify name was updated
	bookmark, err := c.store.GetOne("edited-name")
	if err != nil {
		t.Errorf("Failed to get edited bookmark: %v", err)
	}
	
	if bookmark.Name != "edited-name" {
		t.Errorf("Expected bookmark name 'edited-name', got '%s'", bookmark.Name)
	}
	
	// Verify directory is unchanged
	if bookmark.Directory != testDir1 {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir1, bookmark.Directory)
	}
	
	// Verify old name is gone
	_, err = c.store.GetOne("original-name")
	if err == nil {
		t.Error("Original bookmark name still exists after editing")
	}
	
	// Test: Edit directory only
	err = c.store.Edit("edited-name", "", testDir2)
	if err != nil {
		t.Errorf("Failed to edit bookmark directory: %v", err)
	}
	
	// Verify directory was updated
	bookmark, err = c.store.GetOne("edited-name")
	if err != nil {
		t.Errorf("Failed to get edited bookmark: %v", err)
	}
	
	if bookmark.Directory != testDir2 {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir2, bookmark.Directory)
	}
	
	// Test: Edit both name and directory
	err = c.store.Edit("edited-name", "final-name", testDir1)
	if err != nil {
		t.Errorf("Failed to edit bookmark both name and directory: %v", err)
	}
	
	// Verify both were updated
	bookmark, err = c.store.GetOne("final-name")
	if err != nil {
		t.Errorf("Failed to get final bookmark: %v", err)
	}
	
	if bookmark.Name != "final-name" {
		t.Errorf("Expected bookmark name 'final-name', got '%s'", bookmark.Name)
	}
	
	if bookmark.Directory != testDir1 {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir1, bookmark.Directory)
	}
}

func TestCLI_RunEditCommand_ErrorCases(t *testing.T) {
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
	
	// Test: Edit nonexistent bookmark
	err = c.store.Edit("nonexistent", "new-name", "/tmp")
	if err == nil {
		t.Error("Expected error editing nonexistent bookmark, but got none")
	}
	
	// Add a test bookmark
	testDir := filepath.Join(tmpDir, "test")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	err = c.store.Add("test-edit", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}
	
	// Test: Edit with invalid new name
	err = c.store.Edit("test-edit", "invalid name", testDir)
	if err == nil {
		t.Error("Expected error editing with invalid name, but got none")
	}
	
	// Test: Edit with invalid directory
	err = c.store.Edit("test-edit", "test-edit", "/nonexistent")
	if err == nil {
		t.Error("Expected error editing with invalid directory, but got none")
	}
	
	// Test: Edit with duplicate name (add another bookmark first)
	testDir2 := filepath.Join(tmpDir, "test2")
	err = os.MkdirAll(testDir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory 2: %v", err)
	}
	
	err = c.store.Add("duplicate", testDir2)
	if err != nil {
		t.Errorf("Failed to add duplicate bookmark: %v", err)
	}
	
	// Now try to edit to the duplicate name
	err = c.store.Edit("test-edit", "duplicate", testDir)
	if err == nil {
		t.Error("Expected error editing to duplicate name, but got none")
	}
}

func TestCLI_RunEditCommand_SameName(t *testing.T) {
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
	
	// Add a test bookmark
	testDir := filepath.Join(tmpDir, "test")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	err = c.store.Add("same-name", testDir)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}
	
	// Test: Edit with same name (should work and not fail)
	err = c.store.Edit("same-name", "same-name", testDir)
	if err != nil {
		t.Errorf("Unexpected error editing with same name: %v", err)
	}
	
	// Verify bookmark still exists
	bookmark, err := c.store.GetOne("same-name")
	if err != nil {
		t.Errorf("Failed to get bookmark after edit with same name: %v", err)
	}
	
	if bookmark.Name != "same-name" {
		t.Errorf("Expected bookmark name 'same-name', got '%s'", bookmark.Name)
	}
	
	if bookmark.Directory != testDir {
		t.Errorf("Expected bookmark directory '%s', got '%s'", testDir, bookmark.Directory)
	}
}