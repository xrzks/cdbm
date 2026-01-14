package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func (c *CLI) RunCdCommand(ctx context.Context, cmd *cli.Command) error {
	firstArg := cmd.Args().Get(0)
	if firstArg == "" {
		cli.ShowRootCommandHelpAndExit(cmd, 1)
		return nil
	}
	bookmark, err := c.store.GetOne(firstArg)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(bookmark.Directory)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	cleanedPath := filepath.Clean(abs)

	fileInfo, err := os.Lstat(cleanedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("bookmark directory no longer exists")
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("security violation: path is a symlink")
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	fmt.Printf("cd %s", shellQuote(cleanedPath))
	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
