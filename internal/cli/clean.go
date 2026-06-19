package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewCleanCommand() *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "remove bookmarks pointing to non-existent directories or with invalid names",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"n"},
				Usage:   "show what would be removed without actually deleting",
			},
			&cli.BoolFlag{
				Name:  "invalid-names",
				Usage: "remove bookmarks with invalid names",
			},
			&cli.BoolFlag{
				Name:  "missing-dirs",
				Usage: "remove bookmarks pointing to non-existent directories",
			},
		},
		Action: c.RunCleanCommand,
	}
}

func (c *CLI) RunCleanCommand(ctx context.Context, cmd *cli.Command) error {
	cleanInvalidNames := cmd.Bool("invalid-names")
	cleanMissingDirs := cmd.Bool("missing-dirs")

	if !cleanInvalidNames && !cleanMissingDirs {
		cleanInvalidNames = true
		cleanMissingDirs = true
	}

	bookmarks := c.store.GetAll()

	if len(bookmarks) == 0 {
		fmt.Println("No bookmarks to clean")
		return nil
	}

	var removed []string
	var skipped []struct {
		name string
		err  error
	}
	for _, bm := range bookmarks {
		var removeReason string
		if cleanInvalidNames {
			if err := c.store.ValidateBookmarkName(bm.Name); err != nil {
				removeReason = fmt.Sprintf("invalid name: %s", err)
			}
		}

		if removeReason == "" && cleanMissingDirs {
			if _, err := os.Lstat(bm.Directory); err != nil {
				if os.IsNotExist(err) {
					removeReason = "directory does not exist"
				} else {
					skipped = append(skipped, struct {
						name string
						err  error
					}{name: bm.Name, err: err})
					continue
				}
			}
		}

		if removeReason != "" {
			if cmd.Bool("dry-run") {
				removed = append(removed, bm.Name)
				c.logDebug("clean", map[string]any{
					"name":      bm.Name,
					"directory": bm.Directory,
					"reason":    removeReason,
					"dry_run":   true,
				})
				continue
			}
			if err := c.store.Delete(bm.Name); err != nil {
				return fmt.Errorf("failed to remove bookmark '%s': %w", bm.Name, err)
			}
			removed = append(removed, bm.Name)
			c.logDebug("clean", map[string]any{
				"name":      bm.Name,
				"directory": bm.Directory,
				"reason":    removeReason,
			})
		}
	}

	if len(removed) == 0 {
		fmt.Println("All bookmarks point to valid directories")
		return nil
	}

	fmt.Printf("Removed %d bookmark(s):\n", len(removed))
	for _, name := range removed {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}
