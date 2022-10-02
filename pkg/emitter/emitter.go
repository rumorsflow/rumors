package emitter

import (
	"context"
	"github.com/olebedev/emitter"
	"github.com/rs/zerolog/log"
)

type Emitter interface {
	Use(pattern string, middlewares ...func(event *emitter.Event))
	On(topic string, middlewares ...func(event *emitter.Event)) <-chan emitter.Event
	Once(topic string, middlewares ...func(event *emitter.Event)) <-chan emitter.Event
	Off(topic string, channels ...<-chan emitter.Event)
	Listeners(topic string) []<-chan emitter.Event
	Topics() []string
	Emit(topic string, args ...any) chan struct{}
	Fire(ctx context.Context, topic string, args ...any)
}

type wrapEmitter struct {
	*emitter.Emitter
}

func NewEmitter(capacity uint) Emitter {
	return &wrapEmitter{
		Emitter: emitter.New(capacity),
	}
}

func (e *wrapEmitter) Fire(ctx context.Context, topic string, args ...any) {
	l := log.Ctx(ctx).With().Str("event", topic).Logger()

	done := e.Emit(topic, args...)

	select {
	case <-done:
		l.Debug().Msg("emit done")
	case <-ctx.Done():
		close(done)
		if ctx.Err() != nil {
			l.Error().Err(ctx.Err()).Msg("context done")
		} else {
			l.Warn().Msg("context done")
		}
	}
}
