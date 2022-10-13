package main

import (
	"github.com/fatih/color"
	"github.com/rumorsflow/rumors/internal/cli"
	"os"
	"path/filepath"
)

var version = "(untracked)"

func main() {
	os.Exit(run())
}

// run this CLI application.
func run() int {
	cmd := cli.NewCommand(filepath.Base(os.Args[0]), version)

	if err := cmd.Execute(); err != nil {
		_, _ = color.New(color.FgHiRed, color.Bold).Fprintln(os.Stderr, err.Error())

		return 1
	}

	return 0
}
