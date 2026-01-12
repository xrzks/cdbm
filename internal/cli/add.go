package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Usage:  "add a new entry",
		Action: c.RunAddCommand,
	}
}

func (c *CLI) RunAddCommand(ctx context.Context, cli *cli.Command) error {
	fmt.Println("Add command executed")
	return nil
}
