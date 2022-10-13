package serve

import (
	"fmt"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/config"
	"github.com/rumorsflow/jobs"
	"github.com/rumorsflow/jobs-client"
	"github.com/rumorsflow/logger"
	"github.com/rumorsflow/mongo"
	"github.com/rumorsflow/redis"
	"github.com/rumorsflow/rumors/internal/container"
	"github.com/rumorsflow/rumors/internal/services/parser"
	"github.com/rumorsflow/rumors/internal/services/room"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/rumorsflow/rumors/internal/tasks/feedimporter"
	"github.com/rumorsflow/rumors/internal/tasks/feeditemgroup"
	"github.com/rumorsflow/rumors/internal/tasks/roombroadcast"
	"github.com/rumorsflow/rumors/internal/tasks/roomstart"
	"github.com/rumorsflow/rumors/internal/tasks/roomupdated"
	"github.com/rumorsflow/rumors/internal/tasks/rumors"
	"github.com/rumorsflow/rumors/internal/tasks/tgupdate"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"github.com/rumorsflow/rumors/internal/tgbotserver"
	"github.com/rumorsflow/scheduler"
	"github.com/rumorsflow/scheduler-mongo-provider"
	"github.com/rumorsflow/telegram-bot-api"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

const prefix = "RUMORS"

// NewCommand creates `serve` command.
func NewCommand(cfgFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start Rumors server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			const op = errors.Op("handle serve command")
			// just to be safe
			if cfgFile == "" {
				return errors.E(op, errors.Str("no configuration file provided"))
			}

			// create endure container config
			containerCfg, err := container.NewConfig(cfgFile, prefix)
			if err != nil {
				return errors.E(op, err)
			}

			cfg := &config.Plugin{
				Path:    cfgFile,
				Prefix:  prefix,
				Timeout: containerCfg.GracePeriod,
				Version: version(cmd),
				Cmd:     name(cmd),
			}

			// create endure container
			endureContainer, err := container.NewContainer(*containerCfg)
			if err != nil {
				return errors.E(op, err)
			}

			// register plugins
			err = endureContainer.RegisterAll(
				new(logger.Plugin),
				new(redis.Plugin),
				new(mongo.Plugin),
				new(tgbotapi.Plugin),
				new(tgbotsender.Plugin),
				new(tgbotserver.Plugin),
				new(jobsclient.Plugin),
				new(jobs.Plugin),
				new(scheduler.Plugin),
				new(smp.Plugin),
				new(storage.Plugin),
				new(parser.Plugin),
				new(room.Plugin),
				new(tgupdate.Plugin),
				new(feedimporter.Plugin),
				new(feeditemgroup.Plugin),
				new(roombroadcast.Plugin),
				new(roomstart.Plugin),
				new(roomupdated.Plugin),
				new(rumors.Plugin),
				cfg,
			)
			if err != nil {
				return errors.E(op, err)
			}

			// init container and all services
			err = endureContainer.Init()
			if err != nil {
				return errors.E(op, err)
			}

			// start serving the graph
			errCh, err := endureContainer.Serve()
			if err != nil {
				return errors.E(op, err)
			}

			oss, stop := make(chan os.Signal, 5), make(chan struct{}, 1)
			signal.Notify(oss, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)

			go func() {
				// first catch - stop the container
				<-oss
				// send signal to stop execution
				stop <- struct{}{}

				// after first hit we are waiting for the second
				// second catch - exit from the process
				<-oss
				fmt.Println("exit forced")
				os.Exit(1)
			}()

			for {
				select {
				case e := <-errCh:
					return fmt.Errorf("error: %w\nplugin: %s", e.Error, e.VertexID)
				case <-stop: // stop the container after first signal
					fmt.Printf("stop signal received, grace timeout is: %0.f seconds\n", containerCfg.GracePeriod.Seconds())

					if err = endureContainer.Stop(); err != nil {
						return fmt.Errorf("error: %w", err)
					}

					return nil
				}
			}
		},
	}
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