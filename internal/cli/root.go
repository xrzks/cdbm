package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/logger"
	"github.com/xrzks/cdbm/internal/store"
)

type CLI struct {
	store  *store.Store
	logger logger.Logger
}

func New(s *store.Store) *cli.Command {
	c := &CLI{store: s}
	return &cli.Command{
		Name:                  "cdbm",
		Usage:                 "Manage directory bookmarks",
		Aliases:               []string{"a"},
		Version:               "0.3.5",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"D"},
				Usage:   "enable debug logging to ~/.local/state/cdbm/logs.jsonl",
			},
		},
		Before: c.setupLogger,
		After:  c.closeLogger,
		Commands: []*cli.Command{
			c.NewAddCommand(),
			c.NewListCommand(),
			c.NewInitCommand(),
			c.NewEditCommand(),
			c.NewRemoveCommand(),
		},
		Action: c.RunCdCommand,
	}
}

func (c *CLI) setupLogger(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	if cmd.Bool("debug") {
		statePath, err := config.GetStatePath()
		if err != nil {
			return ctx, fmt.Errorf("failed to initialize debug logging: %w", err)
		}
		logPath := filepath.Join(statePath, "logs.jsonl")
		lgr, err := logger.NewFileLogger(logPath)
		if err != nil {
			return ctx, fmt.Errorf("failed to create debug logger: %w", err)
		}
		c.logger = lgr
	}
	return ctx, nil
}

func (c *CLI) closeLogger(ctx context.Context, cmd *cli.Command) error {
	if c.logger != nil {
		if fl, ok := c.logger.(*logger.FileLogger); ok {
			return fl.Close()
		}
	}
	return nil
}

func (c *CLI) logDebug(action string, details map[string]any) {
	if c.logger != nil {
		if err := c.logger.Log(action, details); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write debug log entry: %v\n", err)
		}
	}
}
