package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type JobOptionDTO struct {
	Type  entity.JobOptionType `json:"type,omitempty" validate:"required,max=50"`
	Value string               `json:"value,omitempty" validate:"required"`
}

type CreateJobDTO struct {
	CronExpr string         `json:"cron_expr,omitempty" validate:"required,min=9,max=254"`
	Name     string         `json:"name,omitempty" validate:"required,max=254"`
	Payload  string         `json:"payload,omitempty"`
	Options  []JobOptionDTO `json:"opts,omitempty" validate:"omitempty,dive"`
	Enabled  bool           `json:"enabled,omitempty"`
}

func (dto CreateJobDTO) toEntity(id uuid.UUID) *entity.Job {
	var options []entity.JobOption

	if len(dto.Options) > 0 {
		options = make([]entity.JobOption, len(dto.Options))
		for i, opt := range dto.Options {
			options[i] = entity.JobOption{
				Type:  opt.Type,
				Value: opt.Value,
			}
		}
	}

	job := &entity.Job{
		ID:       id,
		CronExpr: dto.CronExpr,
		Name:     dto.Name,
	}

	job.SetPayload(dto.Payload)
	job.SetOptions(options)
	job.SetEnabled(dto.Enabled)

	return job
}

type UpdateJobDTO struct {
	CronExpr string          `json:"cron_expr,omitempty" validate:"omitempty,min=9,max=254"`
	Name     string          `json:"name,omitempty" validate:"omitempty,max=254"`
	Payload  *string         `json:"payload,omitempty"`
	Options  *[]JobOptionDTO `json:"opts,omitempty" validate:"omitempty,dive"`
	Enabled  *bool           `json:"enabled,omitempty"`
}

func (dto UpdateJobDTO) toEntity(id uuid.UUID) *entity.Job {
	job := &entity.Job{
		ID:       id,
		CronExpr: dto.CronExpr,
		Name:     dto.Name,
		Payload:  dto.Payload,
		Enabled:  dto.Enabled,
	}

	if dto.Options != nil {
		var options []entity.JobOption
		options = make([]entity.JobOption, len(*dto.Options))
		for i, opt := range *dto.Options {
			options[i] = entity.JobOption{
				Type:  opt.Type,
				Value: opt.Value,
			}
		}
		job.SetOptions(options)
	}

	return job
}

func NewJobCRUD(
	read repository.ReadRepository[*entity.Job],
	write repository.WriteRepository[*entity.Job],
) action.CRUD {
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
