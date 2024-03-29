package task

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"golang.org/x/exp/slog"
	"strings"
	"unicode/utf8"
)

type HandlerTgCmdRumors struct {
	logger      *slog.Logger
	publisher   common.Pub
	articleRepo repository.ReadRepository[*entity.Article]
}

func (h *HandlerTgCmdRumors) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	message := ctx.Value(ctxMsgKey{}).(tgbotapi.Message)
	sites := ctx.Value(ctxSitesKey{}).([]*entity.Site)

	ids := make([]string, len(sites))
	for i, site := range sites {
		ids[i] = site.ID.String()
	}

	grouped, err := h.articles(ctx, message.CommandArguments(), ids)
	if err != nil {
		h.logger.Error("error due to find articles", "err", err, "command", message.Command(), "args", message.CommandArguments(), "telegram_id", message.Chat.ID)
		return err
	}

	h.publisher.Telegram(ctx, model.Message{
		ChatID: message.Chat.ID,
		View:   model.ViewArticles,
		Data:   grouped,
	})

	return nil
}

func (h *HandlerTgCmdRumors) articles(ctx context.Context, args string, siteIDs []string) (map[string][]model.Article, error) {
	index, size, search := pagination(args)
	query := fmt.Sprintf("sort=-pub_date&field.0.0=site_id&cond.0.0=in&value.0.0=%s", strings.Join(siteIDs, ","))

	if utf8.RuneCountInString(search) > 0 {
		filters := []string{
			"&field.1.0=link&cond.1.0=regex&value.1.0=%[1]s",
			"field.1.1=title&cond.1.1=regex&value.1.1=%[1]s",
			"field.1.2=long_desc&cond.1.2=regex&value.1.2=%[1]s",
			"field.1.3=categories&cond.1.3=regex&value.1.3=%[1]s",
		}
		query += fmt.Sprintf(strings.Join(filters, "&"), search)
	}

	grouped := make(map[string][]model.Article, len(siteIDs))

	criteria := db.BuildCriteria(query).SetIndex(int64(index)).SetSize(int64(size))

	iter, err := h.articleRepo.FindIter(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("%s %w", OpServerProcessTask, err)
	}

	for iter.Next(ctx) {
		article := iter.Entity()

		d := article.Domain()

		if _, ok := grouped[d]; ok {
			grouped[d] = append(grouped[d], model.ArticleFromEntity(article))
		} else {
			grouped[d] = []model.Article{model.ArticleFromEntity(article)}
		}
	}

	if err = iter.Close(context.Background()); err != nil {
		return nil, fmt.Errorf("%s %w", OpServerProcessTask, err)
	}

	return grouped, nil
}
