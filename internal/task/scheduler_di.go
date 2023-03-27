package task

import (
	"context"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
)

const ConfigSchedulerKey = "task.scheduler"

type SchedulerKey struct{}

func GetScheduler(ctx context.Context, c ...di.Container) (*Scheduler, error) {
	return di.Get[*Scheduler](ctx, SchedulerKey{}, c...)
}

func SchedulerActivator() *di.Activator {
	return &di.Activator{
		Key: SchedulerKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*SchedulerConfig](c.Configurer(), ConfigSchedulerKey)
			if err != nil {
				return nil, nil, errs.E(di.OpFactory, err)
			}
			cfg.Init()

			repository, err := db.GetJobRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return NewScheduler(repository, rdbMaker, WithInterval(cfg.SyncInterval)), nil, nil
		}),
	}
}
