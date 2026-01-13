package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/store"
)

func main() {
	path := "store"
	store, err := store.NewStore(path)
	if err != nil {
		log.Fatalf("An error has occurred during startup\n %v", err)
	}
	app := cli.New(store)
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
