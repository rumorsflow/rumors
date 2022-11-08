package serve

import (
	"fmt"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/http"
	"github.com/rumorsflow/http/middleware/headers"
	"github.com/rumorsflow/http/middleware/logging"
	"github.com/rumorsflow/http/middleware/parallel_requests"
	"github.com/rumorsflow/http/middleware/proxy_headers"
	"github.com/rumorsflow/http/middleware/recovery"
	"github.com/rumorsflow/jobs"
	"github.com/rumorsflow/jobs-client"
	"github.com/rumorsflow/logger"
	"github.com/rumorsflow/mongo"
	"github.com/rumorsflow/redis"
	"github.com/rumorsflow/rumors/internal/api/errorhandler"
	"github.com/rumorsflow/rumors/internal/api/middleware/jwt"
	"github.com/rumorsflow/rumors/internal/api/v1/auth"
	"github.com/rumorsflow/rumors/internal/api/v1/feeditems"
	"github.com/rumorsflow/rumors/internal/api/v1/feeds"
	"github.com/rumorsflow/rumors/internal/api/v1/rooms"
	"github.com/rumorsflow/rumors/internal/api/v1/schedulerjobs"
	"github.com/rumorsflow/rumors/internal/container"
	"github.com/rumorsflow/rumors/internal/pkg/cobracmd"
	"github.com/rumorsflow/rumors/internal/services/parser"
	"github.com/rumorsflow/rumors/internal/services/room"
	"github.com/rumorsflow/rumors/internal/services/token"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/rumorsflow/rumors/internal/tasks/feedimporter"
	"github.com/rumorsflow/rumors/internal/tasks/feeditemgroup"
	"github.com/rumorsflow/rumors/internal/tasks/roombroadcast"
	"github.com/rumorsflow/rumors/internal/tasks/roomstart"
	"github.com/rumorsflow/rumors/internal/tasks/roomupdated"
	"github.com/rumorsflow/rumors/internal/tasks/rumors"
	"github.com/rumorsflow/rumors/internal/tasks/sources"
	"github.com/rumorsflow/rumors/internal/tasks/subscribe"
	"github.com/rumorsflow/rumors/internal/tasks/subscribed"
	"github.com/rumorsflow/rumors/internal/tasks/tgupdate"
	"github.com/rumorsflow/rumors/internal/tasks/unsubscribe"
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

			ioc, err := container.New(cobracmd.FullName(cmd), cmd.Version, cfgFile)
			if err != nil {
				return errors.E(op, err)
			}

			err = ioc.RegisterAll(
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
				new(sources.Plugin),
				new(subscribed.Plugin),
				new(subscribe.Plugin),
				new(unsubscribe.Plugin),
				new(token.Plugin),
				new(http.Plugin),
				new(errorhandler.Plugin),
				new(parallel_requests.Plugin),
				new(proxy_headers.Plugin),
				new(logging.Plugin),
				new(headers.Plugin),
				new(recovery.Plugin),
				new(jwt.Plugin),
				new(auth.Plugin),
				new(feeds.Plugin),
				new(feeditems.Plugin),
				new(schedulerjobs.Plugin),
				new(rooms.Plugin),
			)
			if err != nil {
				return errors.E(op, err)
			}

			err = ioc.Init()
			if err != nil {
				return errors.E(op, err)
			}

			errCh, err := ioc.Serve()
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
				case <-stop:
					fmt.Println("stop signal received")

					if err = ioc.Stop(); err != nil {
						return fmt.Errorf("error: %w", err)
					}

					return nil
				}
			}
		},
	}
}
