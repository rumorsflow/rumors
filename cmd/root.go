package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var RootCmd = &cobra.Command{
	Use:   "rumors",
	Short: "Rumors CLI",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Version = version(cmd)
	},
}

func init() {
	viper.SetEnvPrefix("RUMORS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	viper.AddConfigPath(".")

	cobra.OnInitialize(func() {
		viper.SetConfigFile(viper.GetString("config"))

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
		}
	})

	flagSet := RootCmd.PersistentFlags()
	flagSet.StringP("config", "c", "config.yaml", "application config path")

	flagSet.BoolP("debug", "d", false, "application debug mode (default \"false\")")

	flagSet.Bool("log.console", false, "true for console log output")
	flagSet.StringP("log.level", "l", "info", "application log level")

	flagSet.String("mongodb.uri", "", "mongo db uri")

	flagSet.String("redis.network", "tcp", "redis network, ex: tcp, unix")
	flagSet.String("redis.address", ":6379", "redis address")
	flagSet.String("redis.username", "", "redis username")
	flagSet.String("redis.password", "", "redis password")
	flagSet.Int("redis.db", 0, "by default redis offers 16 databases (0..15)")

	_ = viper.BindPFlags(flagSet)
}

func version(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Version
	}
	return version(cmd.Parent())
}

func name(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Name()
	}
	return fmt.Sprintf("%s %s", name(cmd.Parent()), cmd.Name())
}
