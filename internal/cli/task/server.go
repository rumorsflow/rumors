package task

import (
	"context"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/http"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os/signal"
)

func NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),

				mongodb.Activator("mongo"),

				pubsub.PublisherActivator(),

				db.ArticleActivator(),
				db.ChatActivator(),
				db.FeedActivator(),

				task.ServerMuxActivator(),
				task.ServerActivator(),
				task.MetricsActivator(),

				http.WoolActivator(cmd.Version),
				http.ServerActivator(nil, nil),
			); err != nil {
				return err
			}

			if _, err := task.GetMetrics(ctx); err != nil {
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

			w, err := http.GetWool(ctx)
			if err != nil {
				return err
			}

			w.MountHealth()

			prometheus.Mount(w)

			srv, err := http.GetServer(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return taskSrv.Run(ctx, taskMux)
			})

			g.Go(func() error {
				return srv.Start(w)
			})

			g.Go(func() error {
				<-ctx.Done()

				return srv.GracefulShutdown(context.Background())
			})

			return g.Wait()
		},
	}
}
