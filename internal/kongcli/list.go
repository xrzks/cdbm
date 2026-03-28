package kongcli

import (
	"fmt"

	"github.com/xrzks/cdbm/internal/store"
)

type ListCmd struct {
	Logger Logger `kong:"-"`
}

func (c *ListCmd) Run(store *store.Store) error {
	bookmarks := store.GetAll()
	for _, bookmark := range bookmarks {
		fmt.Println(bookmark.Pretty())
	}

	if c.Logger != nil {
		_ = c.Logger.Log("list", map[string]any{
			"count": len(bookmarks),
		})
	}

	return nil
}
