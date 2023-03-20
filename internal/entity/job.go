package entity

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type (
	JobName       string
	JobOptionType string
)

const (
	JobFeed    JobName = "job:feed"
	JobSitemap JobName = "job:sitemap"
)

const (
	MaxRetryOpt  JobOptionType = "max-retry"
	QueueOpt     JobOptionType = "queue"
	TimeoutOpt   JobOptionType = "timeout"
	DeadlineOpt  JobOptionType = "deadline"
	UniqueOpt    JobOptionType = "unique"
	ProcessAtOpt JobOptionType = "process-at"
	ProcessInOpt JobOptionType = "process-in"
	TaskIDOpt    JobOptionType = "task-id"
	RetentionOpt JobOptionType = "retention"
	GroupOpt     JobOptionType = "group"
)

type JobOption struct {
	Type  JobOptionType `json:"type" bson:"type"`
	Value string        `json:"value" bson:"value"`
}

type FeedPayload struct {
	SiteID uuid.UUID `json:"site_id,omitempty" bson:"site_id,omitempty"`
	Link   string    `json:"link,omitempty" bson:"link,omitempty"`
}

type SitemapPayload struct {
	SiteID uuid.UUID `json:"site_id,omitempty" bson:"site_id,omitempty"`
	Link   string    `json:"link,omitempty" bson:"link,omitempty"`
}

type Job struct {
	ID        uuid.UUID    `json:"id,omitempty" bson:"_id"`
	CronExpr  string       `json:"cron_expr,omitempty" bson:"cron_expr,omitempty"`
	Name      JobName      `json:"name,omitempty" bson:"name,omitempty"`
	Payload   any          `json:"payload,omitempty" bson:"payload,omitempty"`
	Options   *[]JobOption `json:"options,omitempty" bson:"options,omitempty"`
	Enabled   *bool        `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (e *Job) UnmarshalBSON(data []byte) error {
	var job struct {
		ID        uuid.UUID    `json:"id,omitempty" bson:"_id"`
		CronExpr  string       `json:"cron_expr,omitempty" bson:"cron_expr,omitempty"`
		Name      JobName      `json:"name,omitempty" bson:"name,omitempty"`
		Payload   bson.Raw     `json:"payload,omitempty" bson:"payload,omitempty"`
		Options   *[]JobOption `json:"options,omitempty" bson:"options,omitempty"`
		Enabled   *bool        `json:"enabled,omitempty" bson:"enabled,omitempty"`
		CreatedAt time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	if err := bson.Unmarshal(data, &job); err != nil {
		return err
	}

	e.ID = job.ID
	e.CronExpr = job.CronExpr
	e.Name = job.Name
	e.Options = job.Options
	e.Enabled = job.Enabled
	e.CreatedAt = job.CreatedAt
	e.UpdatedAt = job.UpdatedAt

	switch job.Name {
	case JobFeed:
		e.Payload = &FeedPayload{}
	case JobSitemap:
		e.Payload = &SitemapPayload{}
	default:
		return nil
	}

	return bson.Unmarshal(job.Payload, e.Payload)
}

func (e *Job) EntityID() uuid.UUID {
	return e.ID
}

func (e *Job) SetOptions(options []JobOption) *Job {
	e.Options = &options
	return e
}

func (e *Job) SetEnabled(enabled bool) *Job {
	e.Enabled = &enabled
	return e
}

func (e *Job) HasOptions() bool {
	return e.Options != nil
}

func (e *Job) Active() bool {
	return e.Enabled != nil && *e.Enabled
}
