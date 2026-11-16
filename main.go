package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/store"
)

// version is set at build time with -ldflags "-X main.version=v1.2.3".
var version string

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	s, err := store.NewStore(cfg.StorePath)
	if err != nil {
		return fmt.Errorf("failed to open bookmark store: %w", err)
	}
	return cli.New(buildVersion(), s).Run(context.Background(), os.Args)
}

// buildVersion prefers the version stamped by -ldflags, then the module or VCS
// version the Go toolchain records in the binary.
func buildVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}
