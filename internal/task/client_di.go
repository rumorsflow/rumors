package task

import (
	"context"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
)

type ClientKey struct{}

func GetClient(ctx context.Context, c ...di.Container) (*Client, error) {
	return di.Get[*Client](ctx, ClientKey{}, c...)
}

func ClientActivator() *di.Activator {
	return &di.Activator{
		Key: ClientKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			tc := NewClient(rdbMaker)

			return tc, di.CloserFunc(func(context.Context) error {
				return tc.Close()
			}), nil
		}),
	}
}
