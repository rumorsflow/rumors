package http

import (
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/swagger"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/http"
	"github.com/rumorsflow/rumors/v2/internal/http/front"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os/signal"
)

func NewFrontCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "front",
		Short: "Start Front API Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				mongodb.Activator("mongo"),

				db.ArticleActivator(),
				db.FeedActivator(),

				http.FrontActivator(cmd.Version),
				http.WoolActivator(cmd.Version),
				http.ServerActivator(nil, nil),
			); err != nil {
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
				w.GET("/swagger/front/...", swagger.New(&swagger.Config{InstanceName: "front"}).Handler)
			}

			frontApi.Register(w)

			srv, err := http.GetServer(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return frontApi.Listen(ctx)
			})

			g.Go(func() error {
				return srv.StartC(ctx, w)
			})

			return g.Wait()
		},
	}
}
