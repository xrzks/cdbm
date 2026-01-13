package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a new entry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "the name of the bookmark",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "directory",
				Usage:    "the bookmarked directory",
				Required: true,
			},
		},
		Action: c.RunAddCommand,
	}
}

func (c *CLI) RunAddCommand(ctx context.Context, cli *cli.Command) error {
	directory := cli.String("directory")
	name := cli.String("name")
	return c.store.Add(name, directory)
}
