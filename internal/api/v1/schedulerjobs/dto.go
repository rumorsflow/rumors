package schedulerjobs

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/scheduler-mongo-provider"
	"github.com/samber/lo"
)

type Option struct {
	Type  smp.OptionType `json:"type,omitempty" validate:"required,max=50"`
	Value string         `json:"value,omitempty" validate:"required"`
}

type CreateRequest struct {
	CronExpr string   `json:"cron_expr,omitempty" validate:"required,min=9,max=254"`
	JobCode  string   `json:"job_code,omitempty" validate:"required,max=254"`
	Payload  string   `json:"payload,omitempty"`
	Opts     []Option `json:"opts,omitempty" validate:"omitempty,dive"`
	Enabled  bool     `json:"enabled,omitempty"`
}

func (r CreateRequest) PeriodicTask() smp.PeriodicTask {
	model := smp.PeriodicTask{
		Id:       uuid.NewString(),
		CronExpr: r.CronExpr,
		JobCode:  r.JobCode,
	}

	if len(r.Opts) > 0 {
		model.SetOpts(lo.Map(r.Opts, func(item Option, _ int) smp.Option {
			return smp.Option{Type: item.Type, Value: item.Value}
		}))
	}

	if r.Payload != "" {
		model.SetPayload(r.Payload)
	}

	model.SetEnabled(r.Enabled)

	return model
}

type UpdateRequest struct {
	CronExpr string    `json:"cron_expr,omitempty" validate:"omitempty,min=9,max=254"`
	JobCode  string    `json:"job_code,omitempty" validate:"omitempty,max=254"`
	Payload  *string   `json:"payload,omitempty"`
	Opts     *[]Option `json:"opts,omitempty" validate:"omitempty,dive"`
	Enabled  *bool     `json:"enabled,omitempty"`
}

func (r UpdateRequest) PeriodicTask(id string) smp.PeriodicTask {
	model := smp.PeriodicTask{
		Id:       id,
		CronExpr: r.CronExpr,
		JobCode:  r.JobCode,
		Payload:  r.Payload,
		Enabled:  r.Enabled,
	}

	if r.Opts != nil {
		model.SetOpts(lo.Map(*r.Opts, func(item Option, _ int) smp.Option {
			return smp.Option{Type: item.Type, Value: item.Value}
		}))
	}
	return model
}
