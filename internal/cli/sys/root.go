package sys

import (
	"github.com/rumorsflow/rumors/v2/internal/cli/sys/user"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "sys"}

	cmd.AddCommand(user.NewRootCommand())

	return cmd
}
