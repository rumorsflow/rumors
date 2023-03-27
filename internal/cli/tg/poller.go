package tg

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
	"github.com/rumorsflow/rumors/v2/internal/telegram/poller"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os/signal"
)

func NewPollerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "poller",
		Short: "Start Poller",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),

				telegram.BotActivator(),

				poller.TelegramPollerActivator(),

				task.ClientActivator(),
			); err != nil {
				return err
			}

			tgPoller, err := poller.GetTelegramPoller(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return tgPoller.Poll(ctx)
			})

			g.Go(func() error {
				<-ctx.Done()

				return nil
			})

			return g.Wait()
		},
	}
}
