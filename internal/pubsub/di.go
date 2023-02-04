package pubsub

import (
	"context"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
)

type (
	PubKey struct{}
	SubKey struct{}
)

func GetPub(ctx context.Context, c ...di.Container) (*Publisher, error) {
	return di.Get[*Publisher](ctx, PubKey{}, c...)
}

func GetSub(ctx context.Context, c ...di.Container) (*Subscriber, error) {
	return di.Get[*Subscriber](ctx, SubKey{}, c...)
}

func PublisherActivator() *di.Activator {
	return &di.Activator{
		Key: PubKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			pub := NewPublisher(rdbMaker)

			return pub, di.CloserFunc(func(context.Context) error {
				return pub.Close()
			}), nil
		}),
	}
}

func SubscriberActivator() *di.Activator {
	return &di.Activator{
		Key: SubKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			rdbMaker, err := rdb.GetMaker(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			sub := NewSubscriber(rdbMaker)

			return sub, di.CloserFunc(func(context.Context) error {
				return sub.Close()
			}), nil
		}),
	}
}
