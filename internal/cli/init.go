package cli

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/urfave/cli/v3"
)

//go:embed shell_integration.sh
var shellIntegration string

func (c *CLI) NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Usage:  "generate shell initialization code",
		Action: runIinitCommand,
	}
}

func runIinitCommand(ctx context.Context, cmd *cli.Command) error {
	shell := cmd.Args().Get(0)
	err := installShellIntegration(shell)
	if err != nil {
		return err
	}
	return nil
}

func installShellIntegration(shell string) error {
	switch shell {
	case "zsh", "bash":
		fmt.Println(shellIntegration)
	case "":
		return fmt.Errorf("no shell specified. Usage: cdbm init <zsh|bash>")
	default:
		return fmt.Errorf("unsupported shell: %s (supported shells: zsh, bash)", shell)
	}
	return nil
}
