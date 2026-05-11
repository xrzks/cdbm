package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// NewDeleteCommand creates and returns the 'delete' CLI command.
// This command allows users to remove existing bookmarks by name.
func (c *CLI) NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete an existing bookmark",
		Action: c.RunDeleteCommand,
	}
}

// RunDeleteCommand executes the 'delete' command to remove a bookmark.
// It requires a bookmark name as an argument and returns an error if not provided.
func (c *CLI) RunDeleteCommand(ctx context.Context, cmd *cli.Command) error {
	name := cmd.Args().Get(0)
	if name == "" {
		return fmt.Errorf("no bookmark name specified. Usage: cdbm delete <bookmark-name>")
	}
	if err := c.store.Delete(name); err != nil {
		return err
	}
	c.logDebug("delete", map[string]any{
		"name": name,
	})
	fmt.Printf("Deleted bookmark '%s'\n", name)
	return nil
}
