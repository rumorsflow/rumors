package tg

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os/signal"
)

func NewSubCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sub",
		Short: "Start Sub",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),
				mongodb.Activator("mongo"),
				db.SiteActivator(),
				db.ChatActivator(),

				telegram.BotActivator(),
				telegram.SubActivator(),

				pubsub.SubscriberActivator(),
			); err != nil {
				return err
			}

			sub, err := telegram.GetSub(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return sub.Run(ctx)
			})

			g.Go(func() error {
				<-ctx.Done()

				return nil
			})

			return g.Wait()
		},
	}
}
