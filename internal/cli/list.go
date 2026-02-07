package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list existing entries",
		Action: c.RunListCommand,
	}
}

func (c *CLI) RunListCommand(ctx context.Context, cmd *cli.Command) error {
	bookmarks := c.store.GetAll()
	for _, bookmark := range bookmarks {
		fmt.Println(bookmark.Pretty())
	}
	return nil
}
