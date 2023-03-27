package task

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"sync"
	"time"
)

type (
	PreEnqueueFunc  func(task *asynq.Task, opts []asynq.Option)
	PostEnqueueFunc func(info *asynq.TaskInfo, err error)

	SchedulerOption func(*Scheduler)

	Scheduler struct {
		sync.RWMutex

		interval time.Duration
		repo     repository.ReadRepository[*entity.Job]
		log      *slog.Logger
		so       *asynq.SchedulerOpts
		s        *asynq.Scheduler
		m        map[uuid.UUID]running
	}

	running struct {
		entryID   string
		updatedAt time.Time
	}
)

func WithInterval(interval time.Duration) SchedulerOption {
	return func(s *Scheduler) {
		s.interval = interval
	}
}

func WithPreEnqueueFunc(fn PreEnqueueFunc) SchedulerOption {
	return func(s *Scheduler) {
		s.so.PreEnqueueFunc = fn
	}
}

func WithPostEnqueueFunc(fn PostEnqueueFunc) SchedulerOption {
	return func(s *Scheduler) {
		s.so.PostEnqueueFunc = fn
	}
}

func NewScheduler(repo repository.ReadRepository[*entity.Job], rdbMaker *rdb.UniversalClientMaker, options ...SchedulerOption) *Scheduler {
	log := logger.WithGroup("task").WithGroup("scheduler")

	s := &Scheduler{
		interval: 5 * time.Minute,
		repo:     repo,
		log:      log,
		so: &asynq.SchedulerOpts{
			Logger:   &asynqLogger{logger: log},
			LogLevel: level(log),
		},
		m: map[uuid.UUID]running{},
	}

	for _, option := range options {
		option(s)
	}

	s.s = asynq.NewScheduler(rdbMaker, s.so)

	return s
}

func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.sync(ctx); err != nil {
		if errs.IsCanceledOrDeadline(err) {
			return nil
		}
		return errs.E(OpSchedulerStart, err)
	}

	if err := s.s.Start(); err != nil {
		return errs.E(OpSchedulerStart, err)
	}

	ticker := time.NewTicker(s.interval)

	defer func() {
		ticker.Stop()
		s.log.Debug("stopping syncer goroutine")
		s.s.Shutdown()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.sync(context.Background()); err != nil {
				s.log.Error("failed to sync", err)
			}
		}
	}
}

func (s *Scheduler) Add(job *entity.Job) error {
	if job == nil {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[job.ID]; ok {
		if err := s.remove(job.ID); err != nil {
			return errs.E(OpSchedulerSync, err)
		}
	}

	if err := s.add(job); err != nil {
		return errs.E(OpSchedulerSync, err)
	}

	return nil
}

func (s *Scheduler) Remove(id uuid.UUID) error {
	s.Lock()
	defer s.Unlock()

	return s.remove(id)
}

func (s *Scheduler) add(job *entity.Job) (err error) {
	id := job.ID

	var payload []byte
	if job.Payload != nil {
		switch p := job.Payload.(type) {
		case *entity.FeedPayload:
			p.JobID = &id
		case *entity.SitemapPayload:
			p.JobID = &id
		}

		payload, err = json.Marshal(job.Payload)
		if err != nil {
			return errs.E(
				OpSchedulerAdd,
				job.ID,
				fmt.Errorf("error due to marshal `%s` payload with expr `%s`: %w", job.Name, job.CronExpr, err),
			)
		}
	}
	task := asynq.NewTask(string(job.Name), payload)

	entryID, err := s.s.Register(job.CronExpr, task, opts(job)...)
	if err != nil {
		return errs.E(
			OpSchedulerAdd,
			id,
			fmt.Errorf("failed to register job `%s` with expr `%s`: %w", job.Name, job.CronExpr, err),
		)
	}

	s.m[id] = running{entryID: entryID, updatedAt: job.UpdatedAt}

	s.log.Info("successfully registered job", "id", id, "cron_expr", job.CronExpr, "job_name", job.Name)

	return nil
}

func (s *Scheduler) remove(id uuid.UUID) error {
	if err := s.s.Unregister(s.m[id].entryID); err != nil {
		return errs.E(OpSchedulerRemove, id, fmt.Errorf("failed to unregister job: %w", err))
	}

	delete(s.m, id)

	s.log.Info("successfully unregistered job", "id", id)

	return nil
}

func (s *Scheduler) sync(ctx context.Context) error {
	s.log.Debug("scheduler sync")

	criteria := db.BuildCriteria("sort=-updated_at&field.0.0=enabled&value.0.0=true")
	iter, err := s.repo.FindIter(ctx, criteria)
	if err != nil {
		return errs.E(OpSchedulerSync, err)
	}

	s.Lock()
	defer s.Unlock()

	var newJobs []*entity.Job
	allIDs := make(map[uuid.UUID]struct{}, len(s.m))

	for iter.Next(ctx) {
		job := iter.Entity()

		if r, found := s.m[job.ID]; !found || job.UpdatedAt.After(r.updatedAt) {
			newJobs = append(newJobs, job)
		} else {
			allIDs[job.ID] = struct{}{}
		}
	}

	if err = iter.Close(context.Background()); err != nil {
		return errs.E(OpSchedulerSync, err)
	}

	failed := make(map[uuid.UUID]struct{})

	for id, _ := range s.m {
		if _, ok := allIDs[id]; !ok {
			if err = s.remove(id); err != nil {
				failed[id] = struct{}{}

				s.log.Warn("failed job remove", errs.E(OpSchedulerSync, err))
			}
		}
	}

	for _, job := range newJobs {
		if _, ok := failed[job.ID]; ok {
			s.log.Warn("failed job sync", "id", job.ID)
		} else if err = s.add(job); err != nil {
			s.log.Warn("failed job sync", "id", job.ID)
		}
	}

	return nil
}

func opts(job *entity.Job) []asynq.Option {
	if job.Options != nil {
		options := make([]asynq.Option, len(*job.Options))
		for i, o := range *job.Options {
			options[i] = asynqOpt(o)
		}
		return options
	}
	return nil
}

func asynqOpt(o entity.JobOption) asynq.Option {
	switch o.Type {
	case entity.QueueOpt:
		return asynq.Queue(o.Value)
	case entity.TimeoutOpt:
		return asynq.Timeout(cast.ToDuration(o.Value))
	case entity.DeadlineOpt:
		return asynq.Deadline(cast.ToTime(o.Value))
	case entity.UniqueOpt:
		return asynq.Unique(cast.ToDuration(o.Value))
	case entity.ProcessAtOpt:
		return asynq.ProcessAt(cast.ToTime(o.Value))
	case entity.ProcessInOpt:
		return asynq.ProcessIn(cast.ToDuration(o.Value))
	case entity.TaskIDOpt:
		return asynq.TaskID(o.Value)
	case entity.RetentionOpt:
		return asynq.Retention(cast.ToDuration(o.Value))
	case entity.GroupOpt:
		return asynq.Group(o.Value)
	}
	return asynq.MaxRetry(cast.ToInt(o.Value))
}
