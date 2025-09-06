// Package main is the entrypoint
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/darkweak/rudy/commands"
	"github.com/spf13/cobra"
)

func trapSignals(ctx context.Context, cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	select {
	case <-sig:
		log.Println("[INFO] SIGINT: Stopping RUDY")
		cancel()
	case <-ctx.Done():
		return
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go trapSignals(ctx, cancel)

	var root cobra.Command

	commands.Prepare(&root)

	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
