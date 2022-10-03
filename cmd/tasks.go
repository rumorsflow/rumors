package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var queueCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Start tasks server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.LocalFlags())
	},
	RunE: queue,
}

func init() {
	flagSet := schedulerCmd.PersistentFlags()

	flagSet.Int("asynq.server.concurrency", 0, "how many concurrent workers to use, zero or negative for number of CPUs")
	flagSet.Int("asynq.server.group.max.size", 50, "if zero no delay limit is used")
	flagSet.Duration("asynq.server.group.max.delay", 10*time.Minute, "if zero no size limit is used")
	flagSet.Duration("asynq.server.group.grace.period", 2*time.Minute, "min 1 second")

	RootCmd.AddCommand(queueCmd)
}

func queue(cmd *cobra.Command, _ []string) error {
	return nil
}
