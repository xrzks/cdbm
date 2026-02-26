package cli

import (
	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/store"
)

type CLI struct {
	store *store.Store
}

func New(s *store.Store) *cli.Command {
	c := &CLI{store: s}
	return &cli.Command{
		Name:                  "cdbm",
		Usage:                 "Manage directory bookmarks",
		Aliases:               []string{"a"},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			c.NewAddCommand(),
			c.NewListCommand(),
			c.NewInitCommand(),
			c.NewEditCommand(),
			c.NewDeleteCommand(),
		},
		Action: c.RunCdCommand,
	}
}
