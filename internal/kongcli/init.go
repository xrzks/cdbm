package kongcli

import (
	_ "embed"
	"fmt"

	"github.com/xrzks/cdbm/internal/store"
)

//go:embed shell_integration.sh
var shellIntegration string

//go:embed shell_integration.fish
var shellIntegrationFish string

type InitCmd struct {
	Logger Logger `kong:"-"`
	Shell  string `arg:"" help:"Shell type (zsh, bash, fish)"`
}

func (c *InitCmd) Run(store *store.Store) error {
	switch c.Shell {
	case "zsh", "bash":
		fmt.Println(shellIntegration)
	case "fish":
		fmt.Println(shellIntegrationFish)
	case "":
		return fmt.Errorf("no shell specified. Usage: cdbm init <zsh|bash|fish>")
	default:
		return fmt.Errorf("unsupported shell: %s (supported shells: zsh, bash, fish)", c.Shell)
	}
	return nil
}
