package handlers

import (
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
)

type MetricHandler interface {
	Process(metric *metrics.Metric)
	Init() MetricHandler
}
