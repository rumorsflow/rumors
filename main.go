package main

import (
	"github.com/iagapie/rumors/cmd"
	"log"
)

var version = "(untracked)"

func main() {
	cmd.RootCmd.Version = version

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
