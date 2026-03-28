package main

import (
	"log"

	"github.com/xrzks/cdbm/internal/config"
	kongcli "github.com/xrzks/cdbm/internal/kongcli"
	"github.com/xrzks/cdbm/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store, err := store.NewStore(cfg.StorePath)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	ctx, _ := kongcli.Parse(store)
	if err := ctx.Run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
