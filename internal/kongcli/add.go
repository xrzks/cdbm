package kongcli

import (
	"fmt"
	"os"
	"path/filepath"

	storepkg "github.com/xrzks/cdbm/internal/store"
)

type AddCmd struct {
	Logger    Logger `kong:"-"`
	Name      string `short:"n" help:"The name of the bookmark"`
	Directory string `short:"d" help:"The bookmarked directory"`
}

func (c *AddCmd) Run(s *storepkg.Store) error {
	var name string
	var directory string
	var err error

	if c.Directory != "" {
		directory = c.Directory
	} else {
		directory, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	if c.Name != "" {
		name = c.Name
	} else {
		name = filepath.Base(directory)
		name = storepkg.BookmarkNameSanitizeRegex.ReplaceAllString(name, "")
		if name == "" {
			return fmt.Errorf("cannot generate valid name from directory: %s", directory)
		}
	}

	if err := s.Add(name, directory); err != nil {
		return err
	}

	if c.Logger != nil {
		_ = c.Logger.Log("add", map[string]any{
			"name":      name,
			"directory": directory,
		})
	}

	fmt.Printf("Added bookmark '%s' -> %s\n", name, directory)
	return nil
}
