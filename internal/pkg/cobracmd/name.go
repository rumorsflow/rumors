package cobracmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func FullName(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Name()
	}
	return fmt.Sprintf("%s %s", FullName(cmd.Parent()), cmd.Name())
}
