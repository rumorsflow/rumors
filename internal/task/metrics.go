package task

import (
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"golang.org/x/exp/slog"
)

type Metrics struct {
	logger    *slog.Logger
	inspector *asynq.Inspector
	collector *metrics.QueueMetricsCollector
}

func NewMetrics(rdbMaker *rdb.UniversalClientMaker) *Metrics {
	return &Metrics{
		logger:    logger.WithGroup("task").WithGroup("metrics"),
		inspector: asynq.NewInspector(rdbMaker),
	}
}

func (m *Metrics) Close() error {
	if err := m.inspector.Close(); err != nil {
		return errs.E(OpMetricsClose, err)
	}
	return nil
}

func (m *Metrics) Register() error {
	m.collector = metrics.NewQueueMetricsCollector(m.inspector)

	if err := prometheus.Register(m.collector); err != nil {
		return errs.E(OpMetricsRegister, err)
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
