package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type StoreInterface interface {
	GetOne(name string) (*Bookmark, error)
	Add(name string, directory string) error
	GetAll() []*Bookmark
	Delete(name string) error
	Edit(name string, newName string, newDirectory string) error
}

var BookmarkNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
var BookmarkNameSanitizeRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

type Store struct {
	bookmarks map[string]*Bookmark
	path      string
	mu        sync.RWMutex
}

func NewStore(path string) (*Store, error) {
	store := Store{
		path:      path,
		bookmarks: make(map[string]*Bookmark),
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create store directory: %w", err)
	}

	err := store.loadBookmarks()
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
	if !BookmarkNameRegex.MatchString(name) {
		return fmt.Errorf("bookmark name contains invalid characters (only letters, numbers, ., _, and - allowed)")
	}
	return nil
}

func validateDirectory(directory string) (string, error) {
	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("invalid directory path: %w", err)
	}

	// Check if the final path component is a symlink for security
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

	return filepath.Clean(absPath), nil
}

func (s *Store) GetOne(name string) (*Bookmark, error) {
	if err := validateBookmarkName(name); err != nil {
		return nil, err
	}

	s.mu.RLock()
	bm, exists := s.bookmarks[name]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("bookmark not found")
	}
	return &Bookmark{Name: bm.Name, Directory: bm.Directory}, nil
}

func (s *Store) Add(name string, directory string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*Bookmark, 0, len(s.bookmarks))
	for _, bm := range s.bookmarks {
		list = append(list, &Bookmark{Name: bm.Name, Directory: bm.Directory})
	}
	return list
}

func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.bookmarks[name]; !exists {
		return fmt.Errorf("bookmark not found")
	}
	delete(s.bookmarks, name)
	return s.writeFile()
}

func (s *Store) Edit(name string, newName string, newDirectory string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bm, exists := s.bookmarks[name]
	if !exists {
		return fmt.Errorf("bookmark not found")
	}

	// Determine final name
	finalName := name
	if newName != "" {
		if newName != name {
			if err := validateBookmarkName(newName); err != nil {
				return err
			}
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

	// Create a new bookmark to avoid mutating shared references
	newBookmark := &Bookmark{
		Name:      finalName,
		Directory: finalDirectory,
	}

	// Remove old entry if name changed
	if finalName != name {
		delete(s.bookmarks, name)
	}

	s.bookmarks[finalName] = newBookmark
	return s.writeFile()
}

func (s *Store) loadBookmarks() error {
	data, err := s.loadFile()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s.bookmarks)
	if err != nil {
		return fmt.Errorf("failed to parse bookmarks. The file may be corrupted: %w", err)
	}

	validBookmarks := make(map[string]*Bookmark)
	for name, bm := range s.bookmarks {
		if err := validateBookmarkName(name); err != nil {
			return fmt.Errorf("invalid bookmark name '%s' in store file: %w", name, err)
		}
		if bm == nil {
			return fmt.Errorf("bookmark '%s' is null in store file", name)
		}
		validBookmarks[name] = bm
	}
	s.bookmarks = validBookmarks
	return nil
}
