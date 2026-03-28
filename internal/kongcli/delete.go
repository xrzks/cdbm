package kongcli

import (
	"fmt"

	"github.com/xrzks/cdbm/internal/store"
)

type DeleteCmd struct {
	Name   string `arg:"" help:"Name of the bookmark to delete"`
	Logger Logger `kong:"-"`
}

func (c *DeleteCmd) Run(store *store.Store) error {
	name := c.Name
	if name == "" {
		return fmt.Errorf("no name specified")
	}

	if err := store.Delete(name); err != nil {
		return err
	}

	if c.Logger != nil {
		c.Logger.Log("delete", map[string]any{"name": name})
	}
	fmt.Printf("Deleted bookmark '%s'\n", name)
	return nil
}
