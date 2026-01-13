package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *CLI) NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "generate shell initialization code",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			shell := cmd.Args().Get(0)
			switch shell {
			case "zsh", "bash":
				fmt.Printf(`cdbm() {
				if [[ "$1" == "add" ]] || [[ "$1" == "list" ]]; then
				    command cdbm "$@"
				else
				    eval "$(command cdbm "$1")"
				fi
			}`)
			case "":
				return fmt.Errorf("no shell specified. Usage: cdbm init <zsh|bash>")
			default:
				return fmt.Errorf("unsupported shell: %s (supported shells: zsh, bash)", shell)
			}
			return nil
		},
	}
}
