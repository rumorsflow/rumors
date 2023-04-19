package pubsub

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"sync"
)

type Subscriber struct {
	mu     sync.Mutex
	client redis.UniversalClient
	subs   []*redis.PubSub
}

func NewSubscriber(rdbMaker *rdb.UniversalClientMaker) *Subscriber {
	return &Subscriber{client: rdbMaker.Make()}
}

func (s *Subscriber) All(ctx context.Context) *redis.PubSub {
	return s.pSubscribe(ctx, ChannelPrefix+"*")
}

func (s *Subscriber) Telegram(ctx context.Context) *redis.PubSub {
	return s.subscribe(ctx, ChannelTg)
}

func (s *Subscriber) Articles(ctx context.Context) *redis.PubSub {
	return s.subscribe(ctx, ChannelArticles)
}

func (s *Subscriber) subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := s.client.Subscribe(ctx, channels...)

	s.subs = append(s.subs, sub)

	return sub
}

func (s *Subscriber) pSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := s.client.PSubscribe(ctx, channels...)

	s.subs = append(s.subs, sub)

	return sub
}

func (s *Subscriber) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subs {
		err = errs.Append(err, sub.Close())
	}

	if err = errs.Append(err, s.client.Close()); err != nil {
		return fmt.Errorf("%s error: %w", OpClose, err)
	}

	return nil
}
