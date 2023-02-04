package task

import (
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os/signal"
)

func NewSchedulerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "scheduler",
		Short: "Start Scheduler",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), wool.StopSignals...)
			defer cancel()

			if err := di.Activators(
				rdb.UniversalClientActivator("redis"),
				rdb.MakerActivator(),

				mongodb.Activator("mongo"),

				db.JobActivator(),

				task.SchedulerActivator(),
			); err != nil {
				return err
			}

			sched, err := task.GetScheduler(ctx)
			if err != nil {
				return err
			}

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return sched.Run(ctx)
			})

			g.Go(func() error {
				<-ctx.Done()

				return nil
			})

			return g.Wait()
		},
	}
}
