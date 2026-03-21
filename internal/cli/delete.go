package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete an existing bookmark",
		Action: c.RunDeleteCommand,
	}
}

func (c *CLI) RunDeleteCommand(ctx context.Context, cmd *cli.Command) error {
	name := cmd.Args().Get(0)
	if name == "" {
		return fmt.Errorf("no name specified")
	}
	if err := c.store.Delete(name); err != nil {
		return err
	}
	c.logDebug("delete", map[string]any{
		"name": name,
	})
	fmt.Printf("Deleted bookmark '%s'\n", name)
	return nil
}
