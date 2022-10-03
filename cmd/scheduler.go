package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Start cron scheduler",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.LocalFlags())
	},
	RunE: cron,
}

func init() {
	flagSet := schedulerCmd.PersistentFlags()

	flagSet.String("asynq.scheduler.feed", "@every 5m", "feed importer cron spec string or can use \"@every <duration>\" to specify the interval")

	RootCmd.AddCommand(schedulerCmd)
}

func cron(cmd *cobra.Command, _ []string) error {
	return nil
}
