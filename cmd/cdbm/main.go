package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config\n %v", err)
	}
	store, err := store.NewStore(cfg.StorePath)
	if err != nil {
		log.Fatalf("An error has occurred during startup\n %v", err)
	}

	app := cli.New(store)
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func BinaryExists(name string) (bool, error) {
	_, err := exec.LookPath(name)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("unexpected lookup error: %w", err)
	}

	return true, nil
}
