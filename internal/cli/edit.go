package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewEditCommand() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "edit an existing bookmark (rename or move)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "newName",
				Aliases: []string{"n"},
				Usage:   "the new name of the bookmark",
			},
			&cli.StringFlag{
				Name:    "newDirectory",
				Aliases: []string{"d"},
				Usage:   "the new bookmarked directory",
			},
		},
		Action: c.RunEditCommand,
	}
}

func (c *CLI) RunEditCommand(ctx context.Context, cmd *cli.Command) error {
	name := cmd.Args().Get(0)
	if name == "" {
		return fmt.Errorf("no name specified")
	}
	newName := cmd.String("newName")
	newDirectory := cmd.String("newDirectory")

	if newName == "" && newDirectory == "" {
		return fmt.Errorf("at least one of --newName or --newDirectory must be specified")
	}

	return c.store.Edit(name, newName, newDirectory)
}
