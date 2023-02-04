package tg

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "telegram"}

	cmd.AddCommand(NewPollerCommand())
	cmd.AddCommand(NewSubCommand())

	return cmd
}
