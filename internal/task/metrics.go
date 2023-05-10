package task

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

type Metrics struct {
	logger    *slog.Logger
	inspector *asynq.Inspector
	collector *metrics.QueueMetricsCollector
}

func NewMetrics(redisConnOpt asynq.RedisConnOpt, logger *slog.Logger) *Metrics {
	return &Metrics{
		logger:    logger,
		inspector: asynq.NewInspector(redisConnOpt),
	}
}

func (m *Metrics) Close() error {
	if err := m.inspector.Close(); err != nil {
		return fmt.Errorf("%s %w", OpMetricsClose, err)
	}
	return nil
}

func (m *Metrics) Register() error {
	m.collector = metrics.NewQueueMetricsCollector(m.inspector)

	if err := prometheus.Register(m.collector); err != nil {
		return fmt.Errorf("%s %w", OpMetricsRegister, err)
	}

	m.logger.Info("metrics registered")

	return nil
}

func (m *Metrics) Unregister() {
	if m.collector != nil {
		prometheus.Unregister(m.collector)

		m.logger.Info("metrics unregistered")
	}
}
