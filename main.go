package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cdbmcli "github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	store, err := store.NewStore(cfg.StorePath)
	if err != nil {
		log.Fatalf("Error initializing bookmark store: %v", err)
	}

	app := cdbmcli.New(store)
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
