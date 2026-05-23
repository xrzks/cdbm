package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Usage:   "remove an existing bookmark",
		Aliases: []string{"rm"},
		Action:  c.RunRemoveCommand,
	}
}

func (c *CLI) RunRemoveCommand(ctx context.Context, cmd *cli.Command) error {
	name := cmd.Args().Get(0)
	if name == "" {
		return fmt.Errorf("no bookmark name specified. Usage: cdbm remove <bookmark-name>")
	}
	if err := c.store.Delete(name); err != nil {
		return err
	}
	c.logDebug("remove", map[string]any{
		"name": name,
	})
	fmt.Printf("Removed bookmark '%s'\n", name)
	return nil
}
