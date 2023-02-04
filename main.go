package main

import (
	"github.com/fatih/color"
	"github.com/rumorsflow/rumors/v2/internal/cli"
	"os"
)

//go:generate swag f --dir internal/http/front
//go:generate swag f --dir internal/http/sys

var version = "(untracked)"

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
