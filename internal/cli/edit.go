package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// NewEditCommand creates and returns the 'edit' CLI command.
// This command allows users to modify existing bookmarks by changing their name and/or directory.
func (c *CLI) NewEditCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "edit an existing bookmark (rename or move)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "the new name of the bookmark",
			},
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Usage:   "the new bookmarked directory",
			},
		},
		Action: c.RunEditCommand,
	}
}

// RunEditCommand executes the 'edit' command to modify an existing bookmark.
// It requires a bookmark name as an argument and optional --name and --directory flags
// to specify what to change. At least one flag must be provided.
func (c *CLI) RunEditCommand(ctx context.Context, cmd *cli.Command) error {
	name := cmd.Args().Get(0)
	if name == "" {
		return fmt.Errorf("no name specified")
	}
	newName := cmd.String("name")
	newDirectory := cmd.String("directory")

	if newName == "" && newDirectory == "" {
		return fmt.Errorf("at least one of --name or --directory must be specified to edit the bookmark")
	}
	
	// Validate new name if provided
	if newName != "" {
		if err := c.store.ValidateBookmarkName(newName); err != nil {
			return fmt.Errorf("invalid new name: %w", err)
		}
	}

	if err := c.store.Edit(name, newName, newDirectory); err != nil {
		return err
	}

	details := map[string]any{
		"old_name": name,
	}
	if newName != "" {
		details["new_name"] = newName
	}
	if newDirectory != "" {
		details["new_directory"] = newDirectory
	}
	c.logDebug("edit", details)

	if newName != "" {
		fmt.Printf("Edited bookmark '%s' -> name changed to '%s'\n", name, newName)
	}
	if newDirectory != "" {
		if newName != "" {
			fmt.Printf("Edited bookmark '%s' -> directory changed to %s\n", newName, newDirectory)
		} else {
			fmt.Printf("Edited bookmark '%s' -> directory changed to %s\n", name, newDirectory)
		}
	}

	return nil
}
