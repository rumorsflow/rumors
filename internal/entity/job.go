package entity

import (
	"github.com/google/uuid"
	"time"
)

type JobOptionType string

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

type Job struct {
	ID        uuid.UUID    `json:"id,omitempty" bson:"_id"`
	CronExpr  string       `json:"cron_expr,omitempty" bson:"cron_expr,omitempty"`
	Name      string       `json:"name,omitempty" bson:"name,omitempty"`
	Payload   *string      `json:"payload,omitempty" bson:"payload,omitempty"`
	Options   *[]JobOption `json:"opts,omitempty" bson:"opts,omitempty"`
	Enabled   *bool        `json:"enabled,omitempty" bson:"enabled,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (j *Job) EntityID() uuid.UUID {
	return j.ID
}

func (j *Job) SetPayload(payload string) *Job {
	j.Payload = &payload
	return j
}

func (j *Job) SetOptions(options []JobOption) *Job {
	j.Options = &options
	return j
}

func (j *Job) SetEnabled(enabled bool) *Job {
	j.Enabled = &enabled
	return j
}

func (j *Job) HasOptions() bool {
	return j.Options != nil
}

func (j *Job) Active() bool {
	return j.Enabled != nil && *j.Enabled
}
