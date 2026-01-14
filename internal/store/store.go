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

func validateDirectory(directory string) (string, error) {
	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("invalid directory path: %w", err)
	}

	// check if its a symlink
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("directory '%v' does not exist", directory)
		}
		return "", fmt.Errorf("failed to access directory: %w", err)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("symlinks are not allowed for security reasons")
	}
	if !fileInfo.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}

	return absPath, nil
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

	absPath, err := validateDirectory(directory)
	if err != nil {
		return err
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

func (s *Store) Delete(name string) error {
	_, err := s.GetOne(name)
	if err != nil {
		return err
	}
	delete(s.bookmarks, name)
	return s.writeFile()
}

func (s *Store) Edit(name string, newName string, newDirectory string) error {
	bm, err := s.GetOne(name)
	if err != nil {
		return err
	}

	// Determine final name
	finalName := name
	if newName != "" {
		if err := validateBookmarkName(newName); err != nil {
			return err
		}
		if newName != name {
			if _, exists := s.bookmarks[newName]; exists {
				return fmt.Errorf("bookmark '%s' already exists. Use a different name", newName)
			}
			finalName = newName
		}
	}

	// Determine final directory
	finalDirectory := bm.Directory
	if newDirectory != "" {
		absPath, err := validateDirectory(newDirectory)
		if err != nil {
			return err
		}
		finalDirectory = absPath
	}

	// Update the bookmark
	bm.Name = finalName
	bm.Directory = finalDirectory

	// If name changed, delete old entry and add new one
	if finalName != name {
		delete(s.bookmarks, name)
		s.bookmarks[finalName] = bm
	}

	return s.writeFile()
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
