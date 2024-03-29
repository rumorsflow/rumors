package task

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

type Client struct {
	inner  *asynq.Client
	logger *slog.Logger
}

func NewClient(redisConnOpt asynq.RedisConnOpt, logger *slog.Logger) *Client {
	return &Client{
		inner:  asynq.NewClient(redisConnOpt),
		logger: logger,
	}
}

func (c *Client) Close() error {
	return c.inner.Close()
}

func (c *Client) EnqueueTgCmd(ctx context.Context, message *tgbotapi.Message, updateID int) {
	name := TelegramCmd + message.Command()
	taskID := asynq.TaskID(fmt.Sprintf("%s:%d", name, updateID))

	if err := c.Enqueue(ctx, name, message, taskID, asynq.Queue("tgcmd")); err != nil {
		c.logger.Error("error due to enqueue chat message command", "err", err, "option_task_id", taskID, "message", message)
	}
}

func (c *Client) EnqueueTgMemberNew(ctx context.Context, member *tgbotapi.Chat, updateID int) {
	name := TelegramChatNew
	taskID := asynq.TaskID(fmt.Sprintf("%s:%d", name, updateID))

	if err := c.Enqueue(ctx, name, member, taskID, asynq.Queue("tgmember")); err != nil {
		c.logger.Error("error due to enqueue new chat member", "err", err, "option_task_id", taskID, "member", member)
	}
}

func (c *Client) EnqueueTgMemberEdit(ctx context.Context, member *tgbotapi.ChatMemberUpdated, updateID int) {
	name := TelegramChatEdit
	taskID := asynq.TaskID(fmt.Sprintf("%s:%d", name, updateID))

	if err := c.Enqueue(ctx, name, member, taskID, asynq.Queue("tgmember")); err != nil {
		c.logger.Error("error due to enqueue edit chat member", "err", err, "option_task_id", taskID, "member", member)
	}
}

func (c *Client) Enqueue(ctx context.Context, name string, data any, opts ...asynq.Option) error {
	if err := c.enqueue(ctx, name, data, opts...); err != nil {
		return fmt.Errorf("%s error: %w", OpClientEnqueue, err)
	}
	return nil
}

func (c *Client) enqueue(ctx context.Context, name string, data any, opts ...asynq.Option) error {
	payload, err := marshal(data)
	if err != nil {
		return err
	}

	if c.logger.Enabled(ctx, slog.LevelDebug) {
		c.logger.Debug("task enqueue", "task", name, "payload", data)
	} else {
		c.logger.Info("task enqueue", "task", name)
	}

	info, err := c.inner.EnqueueContext(ctx, asynq.NewTask(name, payload), opts...)
	if err != nil {
		return err
	}

	c.logger.Debug("task enqueued", "id", info.ID, "queue", info.Queue, "task", info.Type)

	return nil
}
