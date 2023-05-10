package main

import (
	"github.com/fatih/color"
	"github.com/rumorsflow/rumors/v2/internal/cli"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
)

var version = "(untracked)"

func init() {
	_, _ = maxprocs.Set(maxprocs.Min(1), maxprocs.Logger(func(_ string, _ ...any) {}))
}

func main() {
	os.Exit(run())
}

// run this CLI application.
func run() int {
	cmd := cli.NewCommand(os.Args, version)

	if err := cmd.Execute(); err != nil {
		_, _ = color.New(color.FgHiRed, color.Bold).Fprintln(os.Stderr, err.Error())

		return 1
	}

	return 0
}
