package user

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "user"}

	cmd.AddCommand(NewCreateCommand())
	cmd.AddCommand(NewQRCommand())

	return cmd
}
