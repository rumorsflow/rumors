package sys

import (
	"errors"
	"github.com/gowool/wool"
	"github.com/hibiken/asynq"
)

const QNameParam = "qname"

type QueueActions struct {
	inspector *asynq.Inspector
}

func NewQueueActions(redisConnOpt asynq.RedisConnOpt) *QueueActions {
	return &QueueActions{inspector: asynq.NewInspector(redisConnOpt)}
}

func (a *QueueActions) Close() error {
	return a.inspector.Close()
}

func (a *QueueActions) Delete(c wool.Ctx) error {
	qname := c.Req().PathParam(QNameParam)
	if err := a.inspector.DeleteQueue(qname, false); err != nil {
		if errors.Is(err, asynq.ErrQueueNotFound) {
			return wool.NewErrNotFound(err)
		}

		if errors.Is(err, asynq.ErrQueueNotEmpty) {
			return wool.NewErrBadRequest(err)
		}

		return err
	}

	return c.NoContent()
}

func (a *QueueActions) Pause(c wool.Ctx) error {
	qname := c.Req().PathParam(QNameParam)
	if err := a.inspector.PauseQueue(qname); err != nil {
		return err
	}

	return c.NoContent()
}

func (a *QueueActions) Resume(c wool.Ctx) error {
	qname := c.Req().PathParam(QNameParam)
	if err := a.inspector.UnpauseQueue(qname); err != nil {
		return err
	}

	return c.NoContent()
}
