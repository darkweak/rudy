package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/darkweak/rudy/commands"
	"github.com/darkweak/rudy/request"
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
	request.Context = ctx

	var root cobra.Command

	commands.Prepare(&root)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
