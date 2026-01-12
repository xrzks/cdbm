package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/store"
)

type CLI struct {
	store *store.Store
}

func New(s *store.Store) *cli.Command {
	c := &CLI{store: s}
	return &cli.Command{
		CommandNotFound: c.RunNotFoundCommand,
		Name:            "cdbm",
		Usage:           "cdbm",
		Aliases:         []string{"a"},
		Commands: []*cli.Command{
			c.NewAddCommand(),
			c.NewListCommand(),
		},
	}
}

func (c *CLI) RunNotFoundCommand(ctx context.Context, cmd *cli.Command, commandName string) {
	cli.ShowAppHelp(cmd)
	fmt.Printf("Command '%s' not found\n", commandName)
}
