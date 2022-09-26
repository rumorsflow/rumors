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
	flagSet.Bool("log.colored", false, "colored log output")
	flagSet.StringP("log.level", "l", "info", "application log level")

	_ = viper.BindPFlags(flagSet)
}

func version(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Version
	}
	return version(cmd.Parent())
}
