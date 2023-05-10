package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/otiai10/opengraph/v2"
	"github.com/oxffaa/gopher-parse-sitemap"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"github.com/rumorsflow/rumors/v2/pkg/util"
	"golang.org/x/exp/slog"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

type HandlerJobSitemap struct {
	logger      *slog.Logger
	publisher   common.Pub
	siteRepo    repository.ReadRepository[*entity.Site]
	articleRepo repository.ReadWriteRepository[*entity.Article]
}

func (h *HandlerJobSitemap) ProcessTask(ctx context.Context, task *asynq.Task) error {
	if task.Payload() == nil {
		h.logger.Warn("task payload is empty")
		return nil
	}

	var payload entity.SitemapPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		h.logger.Error("error due to unmarshal sitemap payload", "err", err, "payload", task.Payload())
		return nil
	}

	site, err := h.siteRepo.FindByID(ctx, payload.SiteID)
	if err != nil {
		if errors.Is(err, repository.ErrEntityNotFound) {
			h.logger.Error("error due to find site", "err", err, "id", payload.SiteID)
			return nil
		}
		return fmt.Errorf("%s find site %v error: %w", OpServerProcessTask, payload.SiteID, err)
	}

	if payload.Lang == nil || *payload.Lang == "" {
		if len(site.Languages) > 0 {
			payload.Lang = &site.Languages[0]
		} else {
			h.logger.Warn("fallback language not found", "payload", payload)
			return nil
		}
	}

	if payload.MatchLoc != nil && *payload.MatchLoc != "" {
		if err = addRegex(payload.MatchLoc); err != nil {
			return fmt.Errorf("%s site %v -> compile payload regex match location error: %w", OpServerProcessTask, payload.SiteID, err)
		}
	}

	if payload.SearchLoc != nil && *payload.SearchLoc != "" {
		if err = addRegex(payload.SearchLoc); err != nil {
			return fmt.Errorf("%s site %v -> compile payload regex search location error: %w", OpServerProcessTask, payload.SiteID, err)
		}
	}

	if payload.IsIndex() {
		if err = sitemap.ParseIndexFromSite(ctx, payload.Link, func(e sitemap.IndexEntry) error {
			payload.Link = e.GetLocation()
			return h.process(ctx, payload, site)
		}); !errs.IsCanceledOrDeadline(err) && !errors.Is(err, io.EOF) {
			return fmt.Errorf("%s %w", OpServerProcessTask, err)
		}
		return nil
	}

	if err = h.process(ctx, payload, site); !errs.IsCanceledOrDeadline(err) && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%s %w", OpServerProcessTask, err)
	}
	return nil
}

func (h *HandlerJobSitemap) process(ctx context.Context, payload entity.SitemapPayload, site *entity.Site) error {
	if err := sitemap.ParseFromSite(ctx, payload.Link, func(e sitemap.Entry) error {
		if matchByLoc(payload.MatchLoc, e.GetLocation()) {
			if search := searchByLoc(payload.SearchLoc, e.GetLocation()); search != "" {
				if payload.SearchLink != nil && *payload.SearchLink != "" {
					search = fmt.Sprintf(*payload.SearchLink, search)
				}

				if h.articleExists(ctx, site, search) {
					if payload.StoppingOnDup() {
						return io.EOF
					}
					return nil
				}
			}

			err := h.processEntry(ctx, e, site, *payload.Lang)
			if errors.Is(err, io.EOF) && !payload.StoppingOnDup() {
				return nil
			}

			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("%s error: %w", OpServerParseSitemap, err)
	}
	return nil
}

func (h *HandlerJobSitemap) processEntry(ctx context.Context, entry sitemap.Entry, site *entity.Site, fallbackLang string) error {
	og, err := h.parseOpengraphMeta(ctx, entry.GetLocation())
	if err != nil {
		if errs.IsCanceledOrDeadline(err) {
			return err
		}

		h.logger.Error("error due to parse sitemap location", "err", fmt.Errorf("%s %w", OpServerProcessTask, err), "entry", entry)

		return nil
	}

	article := &entity.Article{
		ID:     uuid.New(),
		SiteID: site.ID,
		Source: entity.SitemapSource,
		Link:   entry.GetLocation(),
		Title:  util.StripHTMLTags(og.Title),
	}

	if date := entry.GetLastModified(); date != nil {
		article.PubDate = *date
	}

	if entry.GetNews() != nil {
		article.Title = entry.GetNews().Title
		article.Lang = entry.GetNews().Publication.Language
		if date := entry.GetNews().GetPublicationDate(); date != nil {
			article.PubDate = *date
		}
		keywords := strings.Split(entry.GetNews().Keywords, ",")
		categories := make([]string, 0, len(keywords))
		for _, category := range keywords {
			if category = util.StripHTMLTags(category); category != "" {
				categories = append(categories, category)
			}
		}
		if len(categories) > 0 {
			article.SetCategories(categories)
		}
	}

	if article.Title == "" {
		h.logger.Warn("article title not found", "entry", entry, "og", og)
		return nil
	}

	if article.PubDate.IsZero() {
		article.PubDate = time.Now()
	}

	if desc := util.StripHTMLTags(og.Description); utf8.RuneCountInString(desc) >= minShortDesc && !strings.EqualFold(article.Title, desc) {
		article.SetShortDesc(desc)
	}

	media := toMedia(og)
	if len(og.Image) == 0 {
		for _, i := range entry.GetImages() {
			media = append(media, entity.Media{URL: i.ImageLocation, Type: entity.ImageType, Meta: map[string]any{
				"alt": i.ImageTitle,
			}})
		}
	}
	if len(media) > 0 {
		article.SetMedia(media)
	}

	if article.Lang == "" && article.ShortDesc != nil {
		article.Lang = whatlanggo.DetectLang(article.Title + " " + *article.ShortDesc).Iso6391()
	}

	if !contains(site.Languages, article.Lang) {
		article.Lang = fallbackLang
	}

	return h.saveArticle(ctx, article)
}

func (h *HandlerJobSitemap) articleExists(ctx context.Context, site *entity.Site, search string) bool {
	query := fmt.Sprintf("field.0.0=site_id&value.0.0=%s&field.1.0=link&cond.1.0=like&value.1.0=%s", site.ID, search)
	criteria := db.BuildCriteria(query)
	if n, err := h.articleRepo.Count(ctx, criteria.Filter); err == nil && n > 0 {
		return true
	}
	return false
}

func (h *HandlerJobSitemap) saveArticle(ctx context.Context, article *entity.Article) error {
	if err := h.articleRepo.Save(ctx, article); err != nil {
		if errs.IsCanceledOrDeadline(err) {
			return err
		}

		if errors.Is(err, repository.ErrDuplicateKey) {
			h.logger.Debug("error due to save article, duplicate key", "article", article)

			return io.EOF
		} else {
			h.logger.Error("error due to save article", "err", err, "article", article)
		}

		return nil
	}

	h.logger.Debug("article saved", "article", article)

	h.publisher.Articles(ctx, []model.Article{model.ArticleFromEntity(article)})

	return nil
}

func (h *HandlerJobSitemap) parseOpengraphMeta(ctx context.Context, link string) (*opengraph.OpenGraph, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	og, err := openGraphFetch(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("%s error: %w", OpServerParseArticle, err)
	}

	h.logger.Debug("article link parsed", "article", og)

	return og, nil
}
