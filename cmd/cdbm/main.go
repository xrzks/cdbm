package main

import (
	"context"
	"log"
	"os"

	"github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/store"
)

func main() {
	path := "store"
	store, err := store.NewStore(path)
	if err != nil {
		log.Fatalf("An error has occured during execution\n %v", err)
	}
	app := cli.New(store)
	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
