package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (s *Store) loadFile() ([]byte, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte("{}"), nil
		}
		return nil, fmt.Errorf("failed to read store file: %w", err)
	}
	return data, nil
}

func (s *Store) writeFile() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	bytes, err := json.Marshal(s.bookmarks)
	if err != nil {
		return fmt.Errorf("failed to marshal bookmarks: %w", err)
	}
	err = os.WriteFile(s.path, bytes, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
