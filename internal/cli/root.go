package cli

import (
	"context"
	"fmt"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/swagger"
	"github.com/gowool/wool"
	"github.com/joho/godotenv"
	cliHTTP "github.com/rumorsflow/rumors/v2/internal/cli/http"
	cliSys "github.com/rumorsflow/rumors/v2/internal/cli/sys"
	cliTask "github.com/rumorsflow/rumors/v2/internal/cli/task"
	cliTg "github.com/rumorsflow/rumors/v2/internal/cli/tg"
	"github.com/rumorsflow/rumors/v2/internal/http"
	"github.com/rumorsflow/rumors/v2/internal/http/front"
	"github.com/rumorsflow/rumors/v2/internal/http/sys"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/internal/telegram/poller"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"path/filepath"
)

const (
	envDotenv = "DOTENV_PATH"
	prefix    = "RUMORS"
)

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

			if absPath, err := filepath.Abs(cfgFile); err == nil {
				cfgFile = absPath
			}

			if v, ok := os.LookupEnv(envDotenv); ok {
				dotenv = v
			}

			if dotenv != "" {
				if err := godotenv.Load(dotenv); err != nil {
					return err
				}
			}

			cfg := config.NewConfigurer(cfgFile, prefix)

			logCfg, err := config.UnmarshalKey[*logger.Config](cfg, "logger")
			if err != nil {
				return err
			}
			if logCfg.Attrs == nil {
				logCfg.Attrs = make(map[string]any)
			}
			logCfg.Attrs["version"] = version

			logger.Init(logCfg)

			di.Init(cfg)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			err1 := di.Close(context.Background())
			err2 := logger.Sync()

			return multierr.Append(err1, err2)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),
				mongodb.Activator("mongo"),

				pubsub.PublisherActivator(),
				pubsub.SubscriberActivator(),

				db.ArticleActivator(),
				db.ChatActivator(),
				db.FeedActivator(),
				db.JobActivator(),
				db.SysUserActivator(),

				telegram.BotActivator(),
				telegram.SubActivator(),
				poller.TelegramPollerActivator(),

				task.ClientActivator(),
				task.SchedulerActivator(),
				task.ServerMuxActivator(),
				task.ServerActivator(),
				task.MetricsActivator(),

				jwt.ConfigActivator("http.jwt"),
				jwt.SignerActivator(),

				http.FrontActivator(cmd.Version),
				http.SysActivator(),
				http.WoolActivator(cmd.Version),
				http.ServerActivator(nil, nil),
			); err != nil {
				return err
			}

			if _, err := task.GetMetrics(ctx); err != nil {
				return err
			}

			sysApi, err := di.Get[*sys.Sys](ctx, http.SysKey{})
			if err != nil {
				return err
			}

			frontApi, err := di.Get[*front.Front](ctx, http.FrontKey{})
			if err != nil {
				return err
			}

			w, err := http.GetWool(ctx)
			if err != nil {
				return err
			}

			w.MountHealth()

			prometheus.Mount(w)

			if logger.IsDebug() {
				w.Group("/swagger", func(sw *wool.Wool) {
					sw.GET("/sys/...", swagger.New(&swagger.Config{InstanceName: "sys"}).Handler)
					sw.GET("/front/...", swagger.New(&swagger.Config{InstanceName: "front"}).Handler)
				})
			}

			w.Group("", func(sw *wool.Wool) {
				sysApi.Register(sw)
			})

			w.Group("", func(fw *wool.Wool) {
				frontApi.Register(fw)
			})

			srv, err := http.GetServer(ctx)
			if err != nil {
				return err
			}

			sched, err := task.GetScheduler(ctx)
			if err != nil {
				return err
			}

			taskMux, err := task.GetServerMux(ctx)
			if err != nil {
				return err
			}

			taskSrv, err := task.GetServer(ctx)
			if err != nil {
				return err
			}

			bot, err := telegram.GetBot(ctx)
			if err != nil {
				return err
			}

			sub, err := telegram.GetSub(ctx)
			if err != nil {
				return err
			}

			tgPoller, err := poller.GetTelegramPoller(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return frontApi.Listen(ctx)
			})

			g.Go(func() error {
				return srv.Start(w)
			})

			g.Go(func() error {
				return sched.Run(ctx)
			})

			g.Go(func() error {
				return taskSrv.Run(ctx, taskMux)
			})

			g.Go(func() error {
				return sub.Run(ctx)
			})

			g.Go(func() error {
				return tgPoller.Poll(ctx)
			})

			g.Go(func() error {
				if err1 := bot.Send(telegram.Message{View: telegram.ViewAppStart}); err != nil {
					err2 := srv.Shutdown(context.Background())

					return multierr.Append(err1, err2)
				}

				<-ctx.Done()

				err1 := bot.Send(telegram.Message{View: telegram.ViewAppStop})
				err2 := srv.GracefulShutdown(context.Background())

				return multierr.Append(err1, err2)
			})

			return g.Wait()
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	f.StringVar(&dotenv, "dotenv", "", fmt.Sprintf("dotenv file [$%s]", envDotenv))

	_ = f.Parse(args[1:])

	cmd.AddCommand(cliTg.NewRootCommand())
	cmd.AddCommand(cliTask.NewRootCommand())
	cmd.AddCommand(cliHTTP.NewRootCommand())
	cmd.AddCommand(cliSys.NewRootCommand())

	return cmd
}
