// Package main is the main entrypoint
package main

import (
	"github.com/darkweak/rudy/commands"
	"github.com/spf13/cobra"
)

func main() {
	var root cobra.Command

	commands.Prepare(&root)

	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
