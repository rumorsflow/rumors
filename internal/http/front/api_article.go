package front

import (
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

type ArticleActions struct {
	ArticleRepo repository.ReadRepository[*entity.Article]
	FeedRepo    repository.ReadRepository[*entity.Feed]
}

func (a *ArticleActions) List(c wool.Ctx) error {
	query := c.Req().URL.Query()

	feedsFilter := bson.M{"enabled": true}

	if query.Has("host") {
		feedsFilter["host"] = bson.M{"$in": strings.Split(query.Get("host"), ",")}
	}

	if query.Has("source_id") {
		feedsFilter["_id"] = bson.M{"$in": strings.Split(query.Get("source_id"), ",")}
	}

	feeds, err := a.FeedRepo.Find(c.Req().Context(), &repository.Criteria{Filter: feedsFilter})
	if err != nil {
		return nil
	}

	sources := make([]uuid.UUID, len(feeds))
	for i, f := range feeds {
		sources[i] = f.ID
	}

	articlesFilter := bson.M{"source_id": bson.M{"$in": sources}}

	if query.Has("lang") {
		articlesFilter["lang"] = bson.M{"$in": strings.Split(query.Get("lang"), ",")}
	}

	total, err := a.ArticleRepo.Count(c.Req().Context(), articlesFilter)
	if err != nil {
		return err
	}

	criteria := &repository.Criteria{
		Sort:   bson.D{{Key: "pub_date", Value: -1}},
		Filter: articlesFilter,
	}
	criteria.SetIndex(cast.ToInt64(query.Get(db.QueryIndex)))
	criteria.SetSize(cast.ToInt64(query.Get(db.QuerySize)))

	response := action.ListResponse{
		Total: total,
		Index: *criteria.Index,
		Size:  *criteria.Size,
	}

	if total > 0 {
		data, err := a.ArticleRepo.Find(c.Req().Context(), criteria)
		if err != nil {
			return err
		}

		mapped := make([]pubsub.Article, len(data))
		for i, item := range data {
			mapped[i] = pubsub.FromEntity(item)
		}
		response.Data = mapped
	}

	return c.JSON(http.StatusOK, response)
}
