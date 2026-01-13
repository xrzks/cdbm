package store

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *Store) loadFile() ([]byte, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("store file not found. Run 'cdbm init' to set up the application")
		}
		return nil, fmt.Errorf("failed to read store file: %w", err)
	}
	return data, nil
}

func (s *Store) writeFile() error {
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
