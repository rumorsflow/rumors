package cli

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rumorsflow/rumors/v2/internal"
	cliSys "github.com/rumorsflow/rumors/v2/internal/cli/sys"
	"github.com/rumorsflow/rumors/v2/internal/container"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const (
	envDotenv = "DOTENV_PATH"
	prefix    = "RUMORS"
)

func NewCommand(args []string, version string) *cobra.Command {
	var (
		cfgFile string
		dotenv  string
	)

	cmd := &cobra.Command{
		Use:           filepath.Base(args[0]),
		Short:         "Rumors CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Version = version

			if absPath, err := filepath.Abs(cfgFile); err == nil {
				cfgFile = absPath
			}

			if v, ok := os.LookupEnv(envDotenv); ok {
				dotenv = v
			}

			if _, err := os.Stat(dotenv); err == nil {
				if err = godotenv.Load(dotenv); err != nil {
					return err
				}
			}

			cont := &container.Container{
				EnvPrefix: prefix,
				CfgFile:   cfgFile,
				Version:   version,
			}

			cmd.SetContext(context.WithValue(cmd.Context(), "container", cont))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Context().Value("container").(*container.Container).Run(internal.Plugins()...)
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	f.StringVar(&dotenv, "dotenv", ".env", fmt.Sprintf("dotenv file [$%s]", envDotenv))

	_ = f.Parse(args[1:])

	cmd.AddCommand(cliSys.NewRootCommand())

	return cmd
}
