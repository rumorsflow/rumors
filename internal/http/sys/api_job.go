package sys

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
)

type JobOptionDTO struct {
	Type  entity.JobOptionType `json:"type,omitempty" validate:"required,max=50"`
	Value string               `json:"value,omitempty" validate:"required"`
}

func (dto JobOptionDTO) toEntity() entity.JobOption {
	return entity.JobOption{
		Type:  dto.Type,
		Value: dto.Value,
	}
}

type FeedPayloadDTO struct {
	SiteID string `json:"site_id,omitempty" validate:"required,uuid4"`
	Link   string `json:"link,omitempty" validate:"required,url"`
}

func (dto FeedPayloadDTO) toEntity() *entity.FeedPayload {
	siteID, _ := uuid.Parse(dto.SiteID)

	return &entity.FeedPayload{
		SiteID: siteID,
		Link:   dto.Link,
	}
}

type SitemapPayloadDTO struct {
	SiteID     string  `json:"site_id,omitempty" validate:"required,uuid4"`
	Link       string  `json:"link,omitempty" validate:"required,url"`
	Lang       *string `json:"lang,omitempty" validate:"omitempty,bcp47_language_tag"`
	MatchLoc   *string `json:"match_loc,omitempty" validate:"omitempty,max=500"`
	SearchLoc  *string `json:"search_loc,omitempty" validate:"omitempty,max=500"`
	SearchLink *string `json:"search_link,omitempty" validate:"omitempty,max=500"`
	Index      *bool   `json:"index,omitempty"`
	StopOnDup  *bool   `json:"stop_on_dup,omitempty"`
}

func (dto SitemapPayloadDTO) toEntity() *entity.SitemapPayload {
	siteID, _ := uuid.Parse(dto.SiteID)

	return &entity.SitemapPayload{
		SiteID:     siteID,
		Link:       dto.Link,
		Lang:       dto.Lang,
		MatchLoc:   dto.MatchLoc,
		SearchLoc:  dto.SearchLoc,
		SearchLink: dto.SearchLink,
		Index:      dto.Index,
		StopOnDup:  dto.StopOnDup,
	}
}

type CreateJobDTO struct {
	CronExpr string         `json:"cron_expr,omitempty" validate:"required,min=9,max=254"`
	Name     entity.JobName `json:"name,omitempty" validate:"required,max=254"`
	Payload  any            `json:"payload,omitempty" validate:"required,dive"`
	Options  []JobOptionDTO `json:"options,omitempty" validate:"omitempty,dive"`
	Enabled  bool           `json:"enabled,omitempty"`
}

func (dto *CreateJobDTO) UnmarshalJSON(data []byte) error {
	var i struct {
		CronExpr string          `json:"cron_expr,omitempty"`
		Name     entity.JobName  `json:"name,omitempty"`
		Payload  json.RawMessage `json:"payload,omitempty"`
		Options  []JobOptionDTO  `json:"options,omitempty"`
		Enabled  bool            `json:"enabled,omitempty"`
	}

	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	dto.CronExpr = i.CronExpr
	dto.Name = i.Name
	dto.Options = i.Options
	dto.Enabled = i.Enabled

	switch dto.Name {
	case entity.JobFeed:
		dto.Payload = &FeedPayloadDTO{}
	case entity.JobSitemap:
		dto.Payload = &SitemapPayloadDTO{}
	default:
		return nil
	}

	return json.Unmarshal(i.Payload, &dto.Payload)
}

func (dto *CreateJobDTO) toEntity(id uuid.UUID) *entity.Job {
	var options []entity.JobOption

	if len(dto.Options) > 0 {
		options = make([]entity.JobOption, len(dto.Options))
		for i, opt := range dto.Options {
			options[i] = opt.toEntity()
		}
	}

	job := &entity.Job{
		ID:       id,
		CronExpr: dto.CronExpr,
		Name:     dto.Name,
	}

	if dto.Payload != nil {
		switch dto.Name {
		case entity.JobFeed:
			job.Payload = dto.Payload.(*FeedPayloadDTO).toEntity()
		case entity.JobSitemap:
			job.Payload = dto.Payload.(*SitemapPayloadDTO).toEntity()
		}
	}

	job.SetOptions(options)
	job.SetEnabled(dto.Enabled)

	return job
}

type UpdateJobDTO struct {
	CronExpr string          `json:"cron_expr,omitempty" validate:"omitempty,min=9,max=254"`
	Name     entity.JobName  `json:"name,omitempty" validate:"omitempty,max=254"`
	Payload  any             `json:"payload,omitempty" validate:"omitempty,dive"`
	Options  *[]JobOptionDTO `json:"options,omitempty" validate:"omitempty,dive"`
	Enabled  *bool           `json:"enabled,omitempty"`
}

func (dto *UpdateJobDTO) UnmarshalJSON(data []byte) error {
	var i struct {
		CronExpr string          `json:"cron_expr,omitempty"`
		Name     entity.JobName  `json:"name,omitempty"`
		Payload  json.RawMessage `json:"payload,omitempty"`
		Options  *[]JobOptionDTO `json:"options,omitempty"`
		Enabled  *bool           `json:"enabled,omitempty"`
	}

	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	dto.CronExpr = i.CronExpr
	dto.Name = i.Name
	dto.Options = i.Options
	dto.Enabled = i.Enabled

	switch dto.Name {
	case entity.JobFeed:
		dto.Payload = &FeedPayloadDTO{}
	case entity.JobSitemap:
		dto.Payload = &SitemapPayloadDTO{}
	default:
		return nil
	}

	return json.Unmarshal(i.Payload, &dto.Payload)
}

func (dto *UpdateJobDTO) toEntity(id uuid.UUID) *entity.Job {
	job := &entity.Job{
		ID:       id,
		CronExpr: dto.CronExpr,
		Name:     dto.Name,
		Enabled:  dto.Enabled,
	}

	if dto.Payload != nil {
		switch dto.Name {
		case entity.JobFeed:
			job.Payload = dto.Payload.(*FeedPayloadDTO).toEntity()
		case entity.JobSitemap:
			job.Payload = dto.Payload.(*SitemapPayloadDTO).toEntity()
		}
	}

	if dto.Options != nil {
		var options []entity.JobOption
		options = make([]entity.JobOption, len(*dto.Options))
		for i, opt := range *dto.Options {
			options[i] = opt.toEntity()
		}
		job.SetOptions(options)
	}

	return job
}

func NewJobCRUD(read repository.ReadRepository[*entity.Job], write repository.WriteRepository[*entity.Job]) action.CRUD {
	return action.NewCRUD[*CreateJobDTO, *UpdateJobDTO, *entity.Job, any](
		read,
		write,
		action.NewDTOFactory[*CreateJobDTO](),
		action.NewDTOFactory[*UpdateJobDTO](),
		action.RequestMapperFunc[*CreateJobDTO, *entity.Job](func(id uuid.UUID, dto *CreateJobDTO) (*entity.Job, error) {
			return dto.toEntity(id), nil
		}),
		action.RequestMapperFunc[*UpdateJobDTO, *entity.Job](func(id uuid.UUID, dto *UpdateJobDTO) (*entity.Job, error) {
			return dto.toEntity(id), nil
		}),
		nil,
	)
}
