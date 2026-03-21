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
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "the new name of the bookmark",
			},
			&cli.StringFlag{
				Name:    "directory",
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
	newName := cmd.String("name")
	newDirectory := cmd.String("directory")

	if newName == "" && newDirectory == "" {
		return fmt.Errorf("at least one of --name or --directory must be specified")
	}

	if err := c.store.Edit(name, newName, newDirectory); err != nil {
		return err
	}

	details := map[string]any{
		"old_name": name,
	}
	if newName != "" {
		details["new_name"] = newName
	}
	if newDirectory != "" {
		details["new_directory"] = newDirectory
	}
	c.logDebug("edit", details)

	if newName != "" {
		fmt.Printf("Edited bookmark '%s' -> name changed to '%s'\n", name, newName)
	}
	if newDirectory != "" {
		if newName != "" {
			fmt.Printf("Edited bookmark '%s' -> directory changed to %s\n", newName, newDirectory)
		} else {
			fmt.Printf("Edited bookmark '%s' -> directory changed to %s\n", name, newDirectory)
		}
	}

	return nil
}
