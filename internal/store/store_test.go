package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if store == nil {
		t.Fatal("NewStore() returned nil store")
	}

	if store.path != storePath {
		t.Errorf("NewStore() path = %v, want %v", store.path, storePath)
	}

	if store.bookmarks == nil {
		t.Error("NewStore() bookmarks map should be initialized")
	}
}

func TestNewStoreCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "nested", "dir", "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if store == nil {
		t.Fatal("NewStore() returned nil store")
	}

	if _, err := os.Stat(filepath.Dir(storePath)); os.IsNotExist(err) {
		t.Error("NewStore() should create nested directories")
	}
}

func TestNewStoreLoadExisting(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	existingData := `{"test":{"name":"test","directory":"` + tmpDir + `"}}`
	err := os.WriteFile(storePath, []byte(existingData), 0o600)
	if err != nil {
		t.Fatalf("failed to write existing data: %v", err)
	}

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if len(store.bookmarks) != 1 {
		t.Errorf("NewStore() loaded %d bookmarks, want 1", len(store.bookmarks))
	}

	if _, exists := store.bookmarks["test"]; !exists {
		t.Error("NewStore() should load existing bookmark 'test'")
	}
}

func TestAdd(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}

	if len(store.bookmarks) != 1 {
		t.Errorf("Add() resulted in %d bookmarks, want 1", len(store.bookmarks))
	}

	bm, exists := store.bookmarks["test"]
	if !exists {
		t.Fatal("Add() should create bookmark 'test'")
	}

	if bm.Name != "test" {
		t.Errorf("Add() bookmark name = %v, want 'test'", bm.Name)
	}
}

func TestAddDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Add("test", testDir)
	if err == nil {
		t.Error("Add() should return error for duplicate name")
	}
}

func TestAddInvalidName(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test name", testDir)
	if err == nil {
		t.Error("Add() should return error for invalid name")
	}
}

func TestAddInvalidDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	err = store.Add("test", "/non/existent/path")
	if err == nil {
		t.Error("Add() should return error for invalid directory")
	}
}

func TestGetOne(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	bm, err := store.GetOne("test")
	if err != nil {
		t.Errorf("GetOne() error = %v", err)
	}

	if bm == nil {
		t.Fatal("GetOne() returned nil bookmark")
	}

	if bm.Name != "test" {
		t.Errorf("GetOne() bookmark name = %v, want 'test'", bm.Name)
	}
}

func TestGetOneNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	_, err = store.GetOne("nonexistent")
	if err == nil {
		t.Error("GetOne() should return error for non-existent bookmark")
	}
}

func TestGetOneInvalidName(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	_, err = store.GetOne("test name")
	if err == nil {
		t.Error("GetOne() should return error for invalid name")
	}
}

func TestGetAll(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir1 := t.TempDir()
	testDir2 := t.TempDir()

	err = store.Add("test1", testDir1)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Add("test2", testDir2)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	bookmarks := store.GetAll()
	if len(bookmarks) != 2 {
		t.Errorf("GetAll() returned %d bookmarks, want 2", len(bookmarks))
	}

	for _, bm := range bookmarks {
		if bm.Name != "test1" && bm.Name != "test2" {
			t.Errorf("GetAll() returned unexpected bookmark: %v", bm.Name)
		}
	}
}

func TestGetAllEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	bookmarks := store.GetAll()
	if len(bookmarks) != 0 {
		t.Errorf("GetAll() returned %d bookmarks, want 0", len(bookmarks))
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Delete("test")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	if _, exists := store.bookmarks["test"]; exists {
		t.Error("Delete() should remove bookmark")
	}
}

func TestDeleteNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	err = store.Delete("nonexistent")
	if err == nil {
		t.Error("Delete() should return error for non-existent bookmark")
	}
}

func TestEditName(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("oldname", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("oldname", "newname", "")
	if err != nil {
		t.Errorf("Edit() error = %v", err)
	}

	if _, exists := store.bookmarks["oldname"]; exists {
		t.Error("Edit() should remove old bookmark")
	}

	if _, exists := store.bookmarks["newname"]; !exists {
		t.Error("Edit() should create bookmark with new name")
	}
}

func TestEditDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	newDir := t.TempDir()

	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("test", "", newDir)
	if err != nil {
		t.Errorf("Edit() error = %v", err)
	}

	bm, exists := store.bookmarks["test"]
	if !exists {
		t.Fatal("Edit() should keep bookmark")
	}

	if bm.Directory != newDir {
		t.Errorf("Edit() directory = %v, want %v", bm.Directory, newDir)
	}
}

func TestEditBoth(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	newDir := t.TempDir()

	err = store.Add("oldname", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("oldname", "newname", newDir)
	if err != nil {
		t.Errorf("Edit() error = %v", err)
	}

	if _, exists := store.bookmarks["oldname"]; exists {
		t.Error("Edit() should remove old bookmark")
	}

	bm, exists := store.bookmarks["newname"]
	if !exists {
		t.Fatal("Edit() should create bookmark with new name")
	}

	if bm.Directory != newDir {
		t.Errorf("Edit() directory = %v, want %v", bm.Directory, newDir)
	}
}

func TestEditDuplicateName(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir1 := t.TempDir()
	testDir2 := t.TempDir()

	err = store.Add("test1", testDir1)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Add("test2", testDir2)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("test1", "test2", "")
	if err == nil {
		t.Error("Edit() should return error for duplicate name")
	}
}

func TestEditNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	err = store.Edit("nonexistent", "newname", "")
	if err == nil {
		t.Error("Edit() should return error for non-existent bookmark")
	}
}

func TestEditInvalidName(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("test", "test name", "")
	if err == nil {
		t.Error("Edit() should return error for invalid new name")
	}
}

func TestEditInvalidDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = store.Edit("test", "", "/non/existent/path")
	if err == nil {
		t.Error("Edit() should return error for invalid new directory")
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store, err := NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	testDir := t.TempDir()
	err = store.Add("test", testDir)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	data, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("failed to read store file: %v", err)
	}

	if len(data) == 0 {
		t.Error("writeFile() should write data to file")
	}
}

func TestLoadFileNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")

	store := &Store{
		path:      storePath,
		bookmarks: make(map[string]*Bookmark),
	}

	data, err := store.loadFile()
	if err != nil {
		t.Errorf("loadFile() error = %v", err)
	}

	if string(data) != "{}" {
		t.Errorf("loadFile() returned %v, want {}", string(data))
	}
}
