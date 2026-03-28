package kongcli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xrzks/cdbm/internal/store"
)

type CdCmd struct {
	Logger Logger `kong:"-"`
	Name   string `arg:"" help:"Bookmark name"`
}

func (c *CdCmd) Run(store *store.Store) error {
	if c.Name == "" {
		return fmt.Errorf("bookmark name is required")
	}

	bookmark, err := store.GetOne(c.Name)
	if err != nil {
		return err
	}

	cleanedPath := filepath.Clean(bookmark.Directory)

	fileInfo, err := os.Lstat(cleanedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("bookmark directory no longer exists")
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("security violation: path is a symlink")
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	fmt.Printf("cd %s\n", shellQuote(cleanedPath))

	if c.Logger != nil {
		_ = c.Logger.Log("cd", map[string]any{
			"name":      c.Name,
			"directory": cleanedPath,
		})
	}

	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
