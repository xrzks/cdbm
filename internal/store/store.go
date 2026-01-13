package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Store struct {
	bookmarks map[string]*Bookmark
	path      string
}

func NewStore(path string) (*Store, error) {
	var store Store
	store = Store{
		path:      path,
		bookmarks: make(map[string]*Bookmark),
	}
	err := store.initBookmarks()
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func validateBookmarkName(name string) error {
	if name == "" {
		return fmt.Errorf("bookmark name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("bookmark name too long (max 100 characters)")
	}
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, name)
	if err != nil || !matched {
		return fmt.Errorf("bookmark name contains invalid characters (only letters, numbers, ., _, and - allowed)")
	}
	return nil
}

func (s *Store) GetOne(name string) (*Bookmark, error) {
	if err := validateBookmarkName(name); err != nil {
		return nil, err
	}

	bm, exists := s.bookmarks[name]
	if !exists {
		return nil, fmt.Errorf("bookmark not found")
	}
	return bm, nil
}

func (s *Store) Add(name string, directory string) error {
	if err := validateBookmarkName(name); err != nil {
		return err
	}

	if _, exists := s.bookmarks[name]; exists {
		return fmt.Errorf("bookmark '%s' already exists. Use a different name", name)
	}

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("invalid directory path: %w", err)
	}

	// Use Lstat to check for symlinks without following them (prevents TOCTOU attacks)
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory '%v' does not exist", directory)
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}

	// Check if it's a symlink first (redundant with Lstat but explicit for clarity)
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symlinks are not allowed for security reasons")
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	bm := &Bookmark{
		Name:      name,
		Directory: absPath,
	}
	s.bookmarks[name] = bm
	return s.writeFile()
}

func (s *Store) GetAll() []*Bookmark {
	list := make([]*Bookmark, 0, len(s.bookmarks))
	for _, bm := range s.bookmarks {
		list = append(list, bm)
	}
	return list
}

func (s *Store) initBookmarks() error {
	data, err := s.loadFile()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &s.bookmarks)
	if err != nil {
		return fmt.Errorf("failed to parse bookmarks. The file may be corrupted: %w", err)
	}
	return nil
}
