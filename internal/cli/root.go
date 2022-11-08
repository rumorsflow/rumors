package cli

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/internal/cli/serve"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const envDotenv string = "DOTENV_PATH"

func NewCommand(args []string, version string) *cobra.Command {
	var cfgFile string
	var dotenv string

	cmd := &cobra.Command{
		Use:           filepath.Base(args[0]),
		Short:         "Rumors CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Version = version

			// cfgFile could be defined by user or default `config.yaml`
			// this check added just to be safe
			if cfgFile == "" {
				return errors.Str("no configuration file provided")
			}

			// try to get the absolute path to the configuration
			if absPath, err := filepath.Abs(cfgFile); err == nil {
				cfgFile = absPath // switch config path to the absolute
			}

			if v, ok := os.LookupEnv(envDotenv); ok { // read path to the dotenv file from environment variable
				dotenv = v
			}

			if dotenv != "" {
				if err := godotenv.Load(dotenv); err != nil {
					return err
				}
			}

			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	f.StringVar(&dotenv, "dotenv", "", fmt.Sprintf("dotenv file [$%s]", envDotenv))

	_ = f.Parse(args[1:])

	cmd.AddCommand(serve.NewCommand(cfgFile))

	return cmd
}
