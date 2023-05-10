package front

import (
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
	"time"
)

type ArticleActions struct {
	ArticleRepo repository.ReadRepository[*entity.Article]
	SiteRepo    repository.ReadRepository[*entity.Site]
}

func (a *ArticleActions) List(c wool.Ctx) error {
	query := c.Req().URL.Query()

	sitesFilter := bson.M{"enabled": true}

	var (
		ids     []uuid.UUID
		domains []string
	)

	if query.Has("sites") {
		for _, item := range strings.Split(query.Get("sites"), ",") {
			if siteID, err := uuid.Parse(item); err == nil {
				ids = append(ids, siteID)
			} else {
				domains = append(domains, item)
			}
		}
	}

	if len(ids) > 0 {
		sitesFilter["_id"] = bson.M{"$in": ids}
	}
	if len(domains) > 0 {
		sitesFilter["domain"] = bson.M{"$in": domains}
	}

	sites, err := a.SiteRepo.Find(c.Req().Context(), &repository.Criteria{Filter: sitesFilter})
	if err != nil {
		return nil
	}

	siteIDs := make([]uuid.UUID, len(sites))
	for i, s := range sites {
		siteIDs[i] = s.ID
	}

	articlesFilter := bson.M{"site_id": bson.M{"$in": siteIDs}}

	if query.Has("dt") {
		if t, err := time.Parse(time.RFC3339, query.Get("dt")); err == nil {
			articlesFilter["created_at"] = bson.M{"$lte": t}
		}
	}

	if query.Has("langs") {
		articlesFilter["lang"] = bson.M{"$in": strings.Split(query.Get("langs"), ",")}
	}

	total, err := a.ArticleRepo.Count(c.Req().Context(), articlesFilter)
	if err != nil {
		return err
	}

	criteria := &repository.Criteria{
		Sort:   bson.D{{Key: "created_at", Value: -1}},
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

		mapped := make([]model.Article, len(data))
		for i, item := range data {
			mapped[i] = model.ArticleFromEntity(item)
		}
		response.Data = mapped
	}

	return c.JSON(http.StatusOK, response)
}
