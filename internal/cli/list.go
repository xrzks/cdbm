package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// NewListCommand creates and returns the 'list' CLI command.
// This command displays all existing bookmarks in a formatted way.
func (c *CLI) NewListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list existing entries",
		Action: c.RunListCommand,
	}
}

// RunListCommand executes the 'list' command to display all bookmarks.
// Each bookmark is displayed with a formatted name and directory path.
func (c *CLI) RunListCommand(ctx context.Context, cmd *cli.Command) error {
	bookmarks := c.store.GetAll()
	for _, bookmark := range bookmarks {
		fmt.Println(bookmark.Pretty())
	}
	return nil
}
