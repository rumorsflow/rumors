package http

import (
	"context"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/swagger"
	"github.com/rumorsflow/rumors/v2/internal/http"
	"github.com/rumorsflow/rumors/v2/internal/http/sys"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewSysCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sys",
		Short: "Start Sys API Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),
				mongodb.Activator("mongo"),

				db.ArticleActivator(),
				db.SiteActivator(),
				db.ChatActivator(),
				db.JobActivator(),
				db.SysUserActivator(),

				jwt.ConfigActivator("http.jwt"),
				jwt.SignerActivator(),

				http.SysActivator(),
				http.WoolActivator(cmd.Version),
				http.ServerActivator(nil, nil),
			); err != nil {
				return err
			}

			sysApi, err := di.Get[*sys.Sys](ctx, http.SysKey{})
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
				w.GET("/swagger/sys/...", swagger.New(&swagger.Config{InstanceName: "sys"}).Handler)
			}

			sysApi.Register(w)

			srv, err := http.GetServer(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return sysApi.Listen(ctx)
			})

			g.Go(func() error {
				return srv.Start(w)
			})

			g.Go(func() error {
				<-ctx.Done()

				return srv.GracefulShutdown(context.Background())
			})

			logger.Debug("press Ctrl+C to stop")

			return g.Wait()
		},
	}
}
