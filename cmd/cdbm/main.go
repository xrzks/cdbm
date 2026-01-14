package main

import (
	"context"
	"log"
	"os"

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
