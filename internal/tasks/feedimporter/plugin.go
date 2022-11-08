package feedimporter

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/services/parser"
	"github.com/rumorsflow/rumors/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

const PluginName = consts.TaskFeedImporter

type Plugin struct {
	log             *zap.Logger
	client          *asynq.Client
	feedParser      parser.FeedParser
	feedStorage     storage.FeedStorage
	feedItemStorage storage.FeedItemStorage
}

func (p *Plugin) Init(
	log *zap.Logger,
	client *asynq.Client,
	feedImporter parser.FeedParser,
	feedStorage storage.FeedStorage,
	feedItemStorage storage.FeedItemStorage,
) error {
	p.log = log
	p.client = client
	p.feedParser = feedImporter
	p.feedStorage = feedStorage
	p.feedItemStorage = feedItemStorage
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	if task.Payload() == nil {
		p.log.Error("import payload is empty")
		return nil
	}

	feed, err := p.feedStorage.FindById(ctx, string(task.Payload()))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.log.Error("error due to find feed", zap.Error(err), zap.ByteString("feedId", task.Payload()))
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := p.feedParser.Parse(ctx, feed)
	if err != nil {
		p.log.Warn("error due to import feed", zap.Error(err), zap.String("id", feed.Id))
		return err
	}

	for _, item := range data {
		feedItem := &item
		if err = p.feedItemStorage.Save(ctx, feedItem); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				p.log.Debug("error due to save feed item", zap.Error(err), zap.Any("feedItem", feedItem))
			} else {
				p.log.Warn("error due to save feed item", zap.Error(err), zap.Any("feedItem", feedItem))
			}
			continue
		}

		payload, _ := json.Marshal(feedItem)

		t := asynq.NewTask(consts.TaskFeedItemAggregate, payload)
		q := asynq.Queue(consts.QueueFeedItems)
		g := asynq.Group(consts.TaskFeedItemGroup)

		if _, err = p.client.EnqueueContext(ctx, t, q, g); err != nil {
			p.log.Error(
				"error due to enqueue feed item",
				zap.Error(err),
				zap.String("task", task.Type()),
				zap.ByteString("payload", task.Payload()),
			)
		}
	}

	return nil
}
