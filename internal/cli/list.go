package cli

import (
	"context"

	"charm.land/lipgloss/v2"
	"github.com/urfave/cli/v3"
)

func (c *CLI) NewListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list bookmarks",
		Action: c.RunListCommand,
	}
}

func (c *CLI) RunListCommand(ctx context.Context, cmd *cli.Command) error {
	bookmarks := c.store.GetAll()
	for _, bookmark := range bookmarks {
		lipgloss.Println(bookmark.Pretty())
	}
	return nil
}
