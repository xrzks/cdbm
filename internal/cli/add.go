package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/store"
)

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
