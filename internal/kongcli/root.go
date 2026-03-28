package kongcli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/logger"
	"github.com/xrzks/cdbm/internal/store"
)

type CLI struct {
	Debug bool `help:"Enable debug logging to ~/.local/state/cdbm/logs.jsonl" short:"D"`

	Add    AddCmd    `cmd:"" help:"Add a new entry"`
	List   ListCmd   `cmd:"" help:"List existing entries"`
	Init   InitCmd   `cmd:"" help:"Generate shell initialization code"`
	Edit   EditCmd   `cmd:"" help:"Edit an existing bookmark (rename or move)"`
	Delete DeleteCmd `cmd:"" help:"Delete an existing bookmark"`
	Cd     CdCmd     `cmd:"" help:"Change directory to bookmark"`
}

type Logger interface {
	Log(action string, details map[string]any) error
}

func New(store *store.Store) *CLI {
	return &CLI{}
}

func (c *CLI) BeforeApply() error {
	return nil
}

func (c *CLI) bindLogger(lgr Logger) {
	c.Add.Logger = lgr
	c.List.Logger = lgr
	c.Edit.Logger = lgr
	c.Delete.Logger = lgr
	c.Cd.Logger = lgr
	c.Init.Logger = lgr
}

func (c *CLI) closeLogger() error {
	if fl, ok := c.Add.Logger.(*logger.FileLogger); ok {
		return fl.Close()
	}
	return nil
}

func Parse(store *store.Store) (*kong.Context, *CLI) {
	cli := New(store)

	args := os.Args[1:]
	if len(args) > 0 {
		firstArg := args[0]
		isCommand := firstArg == "add" || firstArg == "list" || firstArg == "init" || firstArg == "edit" || firstArg == "delete" || firstArg == "cd" || firstArg == "--help" || firstArg == "-h"
		if !isCommand && !strings.HasPrefix(firstArg, "-") {
			args = append([]string{"cd"}, args...)
		}
	}

	os.Args = append([]string{os.Args[0]}, args...)

	for i, arg := range args {
		if arg == "--debug" || arg == "-D" {
			cli.Debug = true
			statePath, err := config.GetStatePath()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get state path: %v\n", err)
			} else {
				logPath := filepath.Join(statePath, "logs.jsonl")
				lgr, err := logger.NewFileLogger(logPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
				} else {
					cli.bindLogger(lgr)
				}
			}
			args = append(args[:i], args[i+1:]...)
			os.Args = append([]string{os.Args[0]}, args...)
			break
		}
	}

	ctx := kong.Parse(cli, kong.Bind(store))
	if cli.Debug {
		cli.closeLogger()
	}
	return ctx, cli
}
