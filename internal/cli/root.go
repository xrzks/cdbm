// Package cli defines the cdbm command-line interface and its subcommands.
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

// New builds the cdbm command backed by s. version is reported by --version.
func New(version string, s *store.Store) *cli.Command {
	c := &CLI{store: s}
	cmd := &cli.Command{
		Name:                  "cdbm",
		Usage:                 "manage directory bookmarks",
		Version:               version,
		UsageText:             "cdbm [global options] <name>\ncdbm [global options] <command> [command options]",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"D"},
				Usage:   "append debug logs to ~/.local/state/cdbm/logs.jsonl",
			},
		},
		Before: c.setupLogger,
		After:  c.closeLogger,
		Commands: []*cli.Command{
			c.NewAddCommand(),
			c.NewListCommand(),
			c.NewEditCommand(),
			c.NewRemoveCommand(),
			c.NewCleanCommand(),
			c.NewInitCommand(),
		},
		Action: c.RunCdCommand,
	}
	cmd.OnUsageError = onUsageError
	for _, sub := range cmd.Commands {
		sub.OnUsageError = onUsageError
	}
	return cmd
}

// onUsageError returns usage errors instead of printing them with the full
// help; main prints the returned error, so this keeps it to a single line.
func onUsageError(_ context.Context, cmd *cli.Command, err error, _ bool) error {
	return fmt.Errorf("%w (see %s --help)", err, cmd.FullName())
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
			fmt.Fprintf(os.Stderr, "warning: failed to write debug log entry: %v\n", err)
		}
	}
}
