package task

import (
	"context"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
)

type MetricsKey struct{}

func GetMetrics(ctx context.Context, c ...di.Container) (*Metrics, error) {
	return di.Get[*Metrics](ctx, MetricsKey{}, c...)
}

func MetricsActivator() *di.Activator {
	return &di.Activator{
		Key: MetricsKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			metrics := NewMetrics(rdbMaker)
			if err = metrics.Register(); err != nil {
				return nil, nil, errs.E(di.OpFactory, err, metrics.Close())
			}

			return metrics, di.CloserFunc(func(context.Context) error {
				metrics.Unregister()
				return metrics.Close()
			}), nil
		}),
	}
}
