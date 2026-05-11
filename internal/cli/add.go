package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/store"
)

// NewAddCommand creates and returns the 'add' CLI command.
// This command allows users to add new directory bookmarks.
func (c *CLI) NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a new entry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "the name of the bookmark",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "directory",
				Aliases:  []string{"d"},
				Usage:    "the bookmarked directory",
				Required: false,
			},
		},
		Action: c.RunAddCommand,
	}
}

// RunAddCommand executes the 'add' command to create a new bookmark.
// It accepts optional --name and --directory flags, or uses the current directory
// and derives a name from it if not provided.
func (c *CLI) RunAddCommand(ctx context.Context, cmd *cli.Command) error {
	var name string
	var directory string
	var err error

	nameFlag := cmd.String("name")
	directoryFlag := cmd.String("directory")

	// get dir from flag or current working dir
	if cmd.IsSet("directory") {
		directory = directoryFlag
	} else {
		directory, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}
	// get name from flag or current working dir
	if cmd.IsSet("name") {
		name = nameFlag
		// Validate name if provided
		if err := c.store.ValidateBookmarkName(name); err != nil {
			return fmt.Errorf("invalid bookmark name: %w", err)
		}
	} else {
		name = filepath.Base(directory)
		name = store.BookmarkNameSanitizeRegex.ReplaceAllString(name, "")
		if name == "" {
			return fmt.Errorf("cannot generate valid name from directory: %s", directory)
		}
	}
	if err := c.store.Add(name, directory); err != nil {
		return err
	}
	c.logDebug("add", map[string]any{
		"name":      name,
		"directory": directory,
	})
	fmt.Printf("Added bookmark '%s' -> %s\n", name, directory)
	return nil
}
