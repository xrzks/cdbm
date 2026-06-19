package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_RunCleanCommand(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	c := &CLI{store: s}

	existingDir := filepath.Join(tmpDir, "existing")
	err = os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = c.store.Add("valid", existingDir)
	if err != nil {
		t.Errorf("Failed to add valid bookmark: %v", err)
	}

	nonExistentPath := filepath.Join(tmpDir, "does-not-exist")
	bookmarksMap := map[string]*store.Bookmark{
		"valid":   {Name: "valid", Directory: existingDir},
		"invalid": {Name: "invalid", Directory: nonExistentPath},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err = store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to reload store: %v", err)
	}

	c = &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}

	err = c.RunCleanCommand(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	_, err = c.store.GetOne("invalid")
	if err == nil {
		t.Error("Invalid bookmark still exists after cleanup")
	}

	_, err = c.store.GetOne("valid")
	if err != nil {
		t.Errorf("Valid bookmark was removed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark after cleanup, got %d", len(bookmarks))
	}
}

func TestCLI_RunCleanCommand_AllValid(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	c := &CLI{store: s}

	dir1 := filepath.Join(tmpDir, "dir1")
	dir2 := filepath.Join(tmpDir, "dir2")

	err = os.MkdirAll(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = os.MkdirAll(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = c.store.Add("bm1", dir1)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}

	err = c.store.Add("bm2", dir2)
	if err != nil {
		t.Errorf("Failed to add bookmark: %v", err)
	}

	bookmarks := c.store.GetAll()
	initialCount := len(bookmarks)

	err = c.RunCleanCommand(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != initialCount {
		t.Errorf("Expected %d bookmarks after cleanup, got %d", initialCount, len(bookmarks))
	}
}

func TestCLI_RunCleanCommand_AllInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	nonExistentPath1 := filepath.Join(tmpDir, "no1")
	nonExistentPath2 := filepath.Join(tmpDir, "no2")

	bookmarksMap := map[string]*store.Bookmark{
		"bm1": {Name: "bm1", Directory: nonExistentPath1},
		"bm2": {Name: "bm2", Directory: nonExistentPath2},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	c := &CLI{store: s}

	err = c.RunCleanCommand(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks after cleanup, got %d", len(bookmarks))
	}
}

func TestCLI_RunCleanCommand_EmptyStore(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	c := &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 0 {
		t.Errorf("Expected empty store, got %d bookmarks", len(bookmarks))
	}

	err = c.RunCleanCommand(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks after cleanup, got %d", len(bookmarks))
	}
}

func TestCLI_RunCleanCommand_InvalidNames(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	existingDir := filepath.Join(tmpDir, "existing")
	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	bookmarksMap := map[string]*store.Bookmark{
		"valid":          {Name: "valid", Directory: existingDir},
		"invalid space":  {Name: "invalid space", Directory: existingDir},
		"invalid@symbol": {Name: "invalid@symbol", Directory: existingDir},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to reload store: %v", err)
	}

	c := &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}

	err = c.RunCleanCommand(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	_, err = c.store.GetOne("invalid space")
	if err == nil {
		t.Error("Invalid bookmark with space in name still exists after cleanup")
	}

	_, err = c.store.GetOne("invalid@symbol")
	if err == nil {
		t.Error("Invalid bookmark with @ symbol in name still exists after cleanup")
	}

	_, err = c.store.GetOne("valid")
	if err != nil {
		t.Errorf("Valid bookmark was removed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark after cleanup, got %d", len(bookmarks))
	}
}

func TestCLI_RunCleanCommand_InvalidNamesOnly(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	existingDir := filepath.Join(tmpDir, "existing")
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist")

	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	bookmarksMap := map[string]*store.Bookmark{
		"valid":          {Name: "valid", Directory: existingDir},
		"invalid space":  {Name: "invalid space", Directory: existingDir},
		"invalid@symbol": {Name: "invalid@symbol", Directory: nonExistentPath},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to reload store: %v", err)
	}

	c := &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}

	testCmd := &cli.Command{}
	testCmd.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "invalid-names"},
		&cli.BoolFlag{Name: "missing-dirs"},
	}
	testCmd.Set("invalid-names", "true")

	err = c.RunCleanCommand(context.Background(), testCmd)
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	_, err = c.store.GetOne("valid")
	if err != nil {
		t.Errorf("Valid bookmark was removed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark after cleanup (only valid), got %d", len(bookmarks))
	}

	if bookmarks[0].Name != "valid" {
		t.Errorf("Expected only valid bookmark to remain, got %s", bookmarks[0].Name)
	}
}

func TestCLI_RunCleanCommand_MissingDirsOnly(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	existingDir := filepath.Join(tmpDir, "existing")
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist")

	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	bookmarksMap := map[string]*store.Bookmark{
		"valid":          {Name: "valid", Directory: existingDir},
		"invalid space":  {Name: "invalid space", Directory: existingDir},
		"invalid@symbol": {Name: "invalid@symbol", Directory: nonExistentPath},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to reload store: %v", err)
	}

	c := &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}

	testCmd := &cli.Command{}
	testCmd.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "invalid-names"},
		&cli.BoolFlag{Name: "missing-dirs"},
	}
	testCmd.Set("missing-dirs", "true")

	err = c.RunCleanCommand(context.Background(), testCmd)
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	_, err = c.store.GetOne("valid")
	if err != nil {
		t.Errorf("Valid bookmark was removed: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks after cleanup (valid and invalid name), got %d", len(bookmarks))
	}

	foundInvalidSpace := false
	for _, bm := range bookmarks {
		if bm.Name == "invalid space" {
			foundInvalidSpace = true
			break
		}
	}
	if !foundInvalidSpace {
		t.Error("Bookmark with invalid name but valid directory should have been kept")
	}
}

func TestCLI_RunCleanCommand_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	existingDir := filepath.Join(tmpDir, "existing")
	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	bookmarksMap := map[string]*store.Bookmark{
		"valid":        {Name: "valid", Directory: existingDir},
		"invalid":      {Name: "invalid", Directory: filepath.Join(tmpDir, "does-not-exist")},
		"invalid-name": {Name: "invalid@symbol", Directory: existingDir},
	}

	data, err := json.Marshal(bookmarksMap)
	if err != nil {
		t.Fatalf("Failed to marshal bookmarks: %v", err)
	}

	err = os.WriteFile(storePath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write store file: %v", err)
	}

	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	c := &CLI{store: s}

	bookmarks := c.store.GetAll()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}

	testCmd := &cli.Command{}
	testCmd.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "dry-run"},
	}
	testCmd.Set("dry-run", "true")

	err = c.RunCleanCommand(context.Background(), testCmd)
	if err != nil {
		t.Errorf("RunCleanCommand failed: %v", err)
	}

	_, err = c.store.GetOne("invalid")
	if err != nil {
		t.Error("Bookmark should still exist after dry-run")
	}

	_, err = c.store.GetOne("invalid-name")
	if err != nil {
		t.Error("Bookmark should still exist after dry-run")
	}

	_, err = c.store.GetOne("valid")
	if err != nil {
		t.Errorf("Valid bookmark should still exist: %v", err)
	}

	bookmarks = c.store.GetAll()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks after dry-run, got %d", len(bookmarks))
	}
}
