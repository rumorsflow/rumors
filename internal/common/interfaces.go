package common

import (
	"context"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rumorsflow/rumors/v2/internal/model"
)

var Success = errors.New("SUCCESS")

type UnitOfWork interface {
	Repository(tp any) (any, error)
}

type RedisMaker interface {
	Make() (redis.UniversalClient, error)
	MakeRedisClient() any
}

type Client interface {
	EnqueueTgCmd(ctx context.Context, message *tgbotapi.Message, updateID int)
	EnqueueTgMemberNew(ctx context.Context, member *tgbotapi.Chat, updateID int)
	EnqueueTgMemberEdit(ctx context.Context, member *tgbotapi.ChatMemberUpdated, updateID int)
	Enqueue(ctx context.Context, name string, data any, opts ...asynq.Option) error
}

type Pub interface {
	Telegram(ctx context.Context, message any)
	Articles(ctx context.Context, articles []model.Article)
}

type Sub interface {
	All(ctx context.Context) *redis.PubSub
	Telegram(ctx context.Context) *redis.PubSub
	Articles(ctx context.Context) *redis.PubSub
}
