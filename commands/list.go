// Package commands handle the commands
package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/darkweak/rudy/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var list = []commandInstanciator{newRun, newServer}

type (
	command interface {
		Info() string
		GetArgs() cobra.PositionalArgs
		GetDescription() string
		GetLongDescription() string
		GetRequiredFlags() []string
		Run() RunCmd
		SetFlags(p *pflag.FlagSet)
	}
	commandInstanciator func() command
	// RunCmd is the command to run.
	RunCmd func(cmd *cobra.Command, args []string)
)

func setVar[T any](flags *pflag.FlagSet, name string, receiver *T, value *T) {
	if value != nil {
		*receiver = *value

		val := any(*value)

		switch val.(type) {
		case string:
			_ = flags.Set(name, val.(string))
		case int64, int32, int16, int8, int:
			_ = flags.Set(name, fmt.Sprint(val))
		case bool:
			_ = flags.Set(name, fmt.Sprint(val))
		case time.Duration:
			_ = flags.Set(name, val.(time.Duration).String())
		}
	}
}

// Prepare is the setup.
func Prepare(root *cobra.Command) {
	_, err := os.Stat(".rudy.yml")

	for _, item := range list {
		var cobraCmd cobra.Command

		instance := item()

		cobraCmd.Use = instance.Info()
		cobraCmd.Short = instance.GetDescription()
		cobraCmd.Long = instance.GetLongDescription()
		cobraCmd.Args = instance.GetArgs()
		cobraCmd.Run = instance.Run()

		flags := cobraCmd.Flags()
		instance.SetFlags(flags)

		if err == nil {
			mfl, err := configuration.NewFileLoader(".rudy.yml")
			if err != nil {
				panic(fmt.Errorf("unable to load .rudy.yml: %v", err))
			}

			setVar(flags, "concurrents", &concurrents, mfl.Concurrents)
			setVar(flags, "filepath", &filepath, mfl.Filepath)
			setVar(flags, "interval", &interval, mfl.Interval)
			setVar(flags, "size", &size, mfl.Size)
			setVar(flags, "tor", &tor, mfl.Tor)
			setVar(flags, "method", &method, mfl.Method)
			setVar(flags, "headers", &headers, mfl.Headers)
			setVar(flags, "url", &url, mfl.URL)
		}

		for _, f := range instance.GetRequiredFlags() {
			_ = cobraCmd.MarkFlagRequired(f)
		}

		root.AddCommand(&cobraCmd)
	}
}
