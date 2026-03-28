package kongcli

import (
	"fmt"

	"github.com/xrzks/cdbm/internal/store"
)

type EditCmd struct {
	Name      string `arg:"" help:"Name of the bookmark to edit" optional:"" name:"bookmark-name"`
	NewName   string `name:"name" short:"n" help:"The new name of the bookmark"`
	Directory string `name:"directory" short:"d" help:"The new bookmarked directory"`
	Logger    Logger `kong:"-"`
}

func (c *EditCmd) Run(store *store.Store) error {
	name := c.Name
	if name == "" {
		return fmt.Errorf("no name specified")
	}
	newName := c.NewName
	newDirectory := c.Directory

	if newName == "" && newDirectory == "" {
		return fmt.Errorf("at least one of --name or --directory must be specified")
	}

	if err := store.Edit(name, newName, newDirectory); err != nil {
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

	if c.Logger != nil {
		_ = c.Logger.Log("edit", details)
	}

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
