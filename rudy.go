package main

import (
	"github.com/darkweak/rudy/commands"
	"github.com/spf13/cobra"
)

func main() {
	var root cobra.Command

	commands.Prepare(&root)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
