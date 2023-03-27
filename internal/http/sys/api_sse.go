package sys

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/wool"
	"github.com/gowool/wool/render"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"golang.org/x/exp/slog"
	"strings"
	"sync"
	"time"
)

const (
	filterQueues           = "queues"
	filterDailyStats       = "daily_stats"
	filterSchedulerEntries = "scheduler_entries"
)

type (
	AuthDTO struct {
		ClientID string `json:"client_id,omitempty" validate:"required,uuid4"`
	}

	QueueInfo struct {
		// Name of the queue.
		Queue string `json:"queue"`
		// Total number of bytes the queue and its tasks require to be stored in redis.
		MemoryUsage int64 `json:"memory_usage_bytes"`
		// Total number of tasks in the queue.
		Size int `json:"size"`
		// Totoal number of groups in the queue.
		Groups int `json:"groups"`
		// Latency of the queue in milliseconds.
		LatencyMillisec int64 `json:"latency_msec"`
		// Latency duration string for display purpose.
		DisplayLatency string `json:"display_latency"`

		// Number of tasks in each state.
		Active      int `json:"active"`
		Pending     int `json:"pending"`
		Aggregating int `json:"aggregating"`
		Scheduled   int `json:"scheduled"`
		Retry       int `json:"retry"`
		Archived    int `json:"archived"`
		Completed   int `json:"completed"`

		// Total number of tasks processed during the given date.
		// The number includes both succeeded and failed tasks.
		Processed int `json:"processed"`
		// Breakdown of processed tasks.
		Succeeded int `json:"succeeded"`
		Failed    int `json:"failed"`
		// Paused indicates whether the queue is paused.
		// If true, tasks in the queue will not be processed.
		Paused bool `json:"paused"`
		// Time when this snapshot was taken.
		Timestamp time.Time `json:"timestamp"`
	}

	DailyStats struct {
		Queue     string `json:"queue"`
		Processed int    `json:"processed"`
		Succeeded int    `json:"succeeded"`
		Failed    int    `json:"failed"`
		Date      string `json:"date"`
	}

	SchedulerEntry struct {
		ID            string          `json:"id"`
		Spec          string          `json:"spec"`
		JobName       string          `json:"job_name"`
		JobPayload    json.RawMessage `json:"job_payload"`
		Opts          []string        `json:"options"`
		NextEnqueueAt time.Time       `json:"next_enqueue_at"`
		// This field is omitted if there were no previous enqueue events.
		PrevEnqueueAt *time.Time `json:"prev_enqueue_at,omitempty"`
	}

	StatsResponse struct {
		Queues           []*QueueInfo             `json:"queues,omitempty"`
		DailyStats       map[string][]*DailyStats `json:"daily_stats,omitempty"`
		SchedulerEntries []*SchedulerEntry        `json:"scheduler_entries,omitempty"`
	}

	SSE struct {
		*sse.Event
		inspector *asynq.Inspector
		logger    *slog.Logger
		clients   sync.Map
	}

	sseClient struct {
		auth             bool
		queues           bool
		dailyStats       bool
		schedulerEntries bool
	}
)

func NewSSE(rdbMaker *rdb.UniversalClientMaker, logger *slog.Logger) *SSE {
	return &SSE{
		Event:     sse.New(&sse.Config{ClientIdle: 5 * time.Minute}, logger),
		inspector: asynq.NewInspector(rdbMaker),
		logger:    logger,
	}
}

func (a *SSE) Auth(c wool.Ctx) error {
	var dto AuthDTO
	if err := c.Bind(&dto); err != nil {
		return err
	}

	if client, ok := a.clients.Load(dto.ClientID); ok {
		client.(*sseClient).auth = true

		return c.NoContent()
	}

	return wool.NewErrBadRequest(nil)
}

func (a *SSE) Middleware(next wool.Handler) wool.Handler {
	return a.Event.Middleware(func(c wool.Ctx) error {
		cl, ok := c.Get(sse.ClientKey).(sse.Client)
		if !ok {
			return errors.New("SSE client not found")
		}

		f := c.Req().QueryParam("filter")

		client := &sseClient{
			queues:           f == "" || strings.Contains(f, filterQueues),
			dailyStats:       f == "" || strings.Contains(f, filterDailyStats),
			schedulerEntries: f == "" || strings.Contains(f, filterSchedulerEntries),
		}

		a.clients.Store(cl.ID, client)
		defer a.clients.Delete(cl.ID)

		cancelCtx, cancelRequest := context.WithCancel(c.Req().Context())
		defer cancelRequest()

		c.SetReq(c.Req().WithContext(cancelCtx))

		t := time.AfterFunc(5*time.Second, func() {
			if auth, ok := a.clients.Load(cl.ID); !ok || auth.(*sseClient).auth == false {
				cancelRequest()
			}
		})
		defer t.Stop()

		return next(c)
	})
}

func (a *SSE) Listen(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)

	defer func() {
		ticker.Stop()
		a.logger.Debug("sse listener stop")
	}()

	a.logger.Debug("sse listener start")

	for {
		select {
		case <-ctx.Done():
			return a.Close()
		case <-ticker.C:
			ticker.Stop()
			a.broadcast()
			ticker.Reset(time.Second)
		}
	}
}

func (a *SSE) broadcast() {
	a.logger.Debug("sse start to accumulate the stats")

	var (
		err              error
		queues           []*QueueInfo
		dailyStats       map[string][]*DailyStats
		schedulerEntries []*SchedulerEntry
	)

	a.clients.Range(func(key, value any) bool {
		client := value.(*sseClient)
		if !client.auth {
			return true
		}

		response := &StatsResponse{}

		if client.queues {
			if queues == nil {
				queues, err = a.queues()
				if err != nil {
					a.logger.Error("error due to collect queues info", err)
					return true
				}
			}
			response.Queues = queues
		}

		if client.dailyStats {
			if dailyStats == nil {
				dailyStats, err = a.dailyStats()
				if err != nil {
					a.logger.Error("error due to collect queues daily stats", err)
					return true
				}
			}
			response.DailyStats = dailyStats
		}

		if client.schedulerEntries {
			if schedulerEntries == nil {
				schedulerEntries, err = a.schedulerEntries()
				if err != nil {
					a.logger.Error("error due to collect scheduler entries", err)
					return true
				}
			}
			response.SchedulerEntries = schedulerEntries
		}

		a.Notify(key.(string), render.SSEvent{
			Event: "stats",
			Data:  response,
		})

		return true
	})

	a.logger.Debug("sse end to accumulate the stats")
}

func (a *SSE) queues() ([]*QueueInfo, error) {
	queues, err := a.inspector.Queues()
	if err != nil {
		return nil, err
	}

	infos := make([]*QueueInfo, len(queues))

	for i, queue := range queues {
		info, err := a.queueInfo(queue)
		if err != nil {
			return nil, err
		}
		infos[i] = info
	}

	return infos, nil
}

func (a *SSE) queueInfo(queue string) (*QueueInfo, error) {
	info, err := a.inspector.GetQueueInfo(queue)
	if err != nil {
		return nil, err
	}

	return &QueueInfo{
		Queue:           info.Queue,
		MemoryUsage:     info.MemoryUsage,
		Size:            info.Size,
		Groups:          info.Groups,
		LatencyMillisec: info.Latency.Milliseconds(),
		DisplayLatency:  info.Latency.Round(10 * time.Millisecond).String(),
		Active:          info.Active,
		Pending:         info.Pending,
		Aggregating:     info.Aggregating,
		Scheduled:       info.Scheduled,
		Retry:           info.Retry,
		Archived:        info.Archived,
		Completed:       info.Completed,
		Processed:       info.Processed,
		Succeeded:       info.Processed - info.Failed,
		Failed:          info.Failed,
		Paused:          info.Paused,
		Timestamp:       info.Timestamp,
	}, nil
}

func (a *SSE) dailyStats() (map[string][]*DailyStats, error) {
	queues, err := a.inspector.Queues()
	if err != nil {
		return nil, err
	}

	dailyStats := make(map[string][]*DailyStats)
	for _, queue := range queues {
		stats, err := a.queueStats(queue)
		if err != nil {
			return nil, err
		}

		dailyStats[queue] = stats
	}

	return dailyStats, nil
}

func (a *SSE) queueStats(queue string) ([]*DailyStats, error) {
	// get stats for the last 90 days
	history, err := a.inspector.History(queue, 90)
	if err != nil {
		return nil, err
	}

	out := make([]*DailyStats, len(history))

	for i, s := range history {
		out[i] = &DailyStats{
			Queue:     s.Queue,
			Processed: s.Processed,
			Succeeded: s.Processed - s.Failed,
			Failed:    s.Failed,
			Date:      s.Date.Format("2006-01-02"),
		}
	}

	return out, nil
}

func (a *SSE) schedulerEntries() ([]*SchedulerEntry, error) {
	entries, err := a.inspector.SchedulerEntries()
	if err != nil {
		return nil, err
	}

	out := make([]*SchedulerEntry, len(entries))
	for i, e := range entries {
		opts := make([]string, len(e.Opts)) // create a non-nil, empty slice to avoid null in json output
		for j, o := range e.Opts {
			opts[j] = o.String()
		}

		var prev *time.Time
		if !e.Prev.IsZero() {
			prev = &e.Prev
		}

		out[i] = &SchedulerEntry{
			ID:            e.ID,
			Spec:          e.Spec,
			JobName:       e.Task.Type(),
			JobPayload:    e.Task.Payload(),
			Opts:          opts,
			NextEnqueueAt: e.Next,
			PrevEnqueueAt: prev,
		}
	}

	return out, nil
}
