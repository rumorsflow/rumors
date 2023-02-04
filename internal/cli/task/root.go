package task

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "task"}

	cmd.AddCommand(NewSchedulerCommand())
	cmd.AddCommand(NewServerCommand())

	return cmd
}
