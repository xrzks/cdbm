package store

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *Store) loadFile() ([]byte, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func (s *Store) writeFile(path string, data []byte) error {
	bytes, err := json.Marshal(s.bookmarks)
	if err != nil {
		return fmt.Errorf("failed to marshal bookmarks: %w", err)
	}
	err = os.WriteFile(s.path, bytes, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
