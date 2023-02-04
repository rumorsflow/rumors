package http

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "http"}

	cmd.AddCommand(NewSysCommand())
	cmd.AddCommand(NewFrontCommand())

	return cmd
}
