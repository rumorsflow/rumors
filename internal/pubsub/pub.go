package pubsub

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"golang.org/x/exp/slog"
)

const (
	ChannelPrefix   = "rumors.event."
	ChannelArticles = ChannelPrefix + "articles"
	ChannelTg       = ChannelPrefix + "telegram"

	OpMarshal = "pubsub: marshal"
	OpPublish = "pubsub: publish"
	OpClose   = "pubsub: close"
)

type Publisher struct {
	client redis.UniversalClient
	logger *slog.Logger
}

func NewPublisher(rdbMaker common.RedisMaker, logger *slog.Logger) (*Publisher, error) {
	client, err := rdbMaker.Make()
	if err != nil {
		return nil, err
	}

	return &Publisher{client: client, logger: logger}, nil
}

func (p *Publisher) Telegram(ctx context.Context, message any) {
	if err := p.publish(ctx, ChannelTg, message); err != nil {
		p.error("error due to publish on telegram", ChannelTg, err)
	}
}

func (p *Publisher) Articles(ctx context.Context, articles []model.Article) {
	if err := p.publish(ctx, ChannelArticles, articles); err != nil {
		p.error("error due to publish articles", ChannelArticles, err)
	}
}

func (p *Publisher) publish(ctx context.Context, channel string, message any) (err error) {
	switch message.(type) {
	case string, []byte:
		break
	default:
		message, err = json.Marshal(message)
		if err != nil {
			return fmt.Errorf("%s %w", OpMarshal, err)
		}
	}

	if err = p.client.Publish(ctx, channel, message).Err(); err != nil {
		return fmt.Errorf("%s %w", OpPublish, err)
	}

	p.logger.Debug("pubsub published a message", "channel", channel, "message", message)

	return nil
}

func (p *Publisher) Close() error {
	if err := p.client.Close(); err != nil {
		return fmt.Errorf("%s %w", OpClose, err)
	}
	return nil
}

func (p *Publisher) error(msg, ch string, err error) {
	p.logger.Error(msg, "err", err, "channel", ch)
}
