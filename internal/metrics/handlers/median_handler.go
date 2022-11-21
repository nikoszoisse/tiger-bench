package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"time"
)

type MedianHandler struct {
	series []time.Duration
}

func (a *MedianHandler) Init() MetricHandler {
	a.series = make([]time.Duration, 0)
	return a
}

func (a *MedianHandler) Process(metric *metrics.Metric) {
	a.series = append(a.series, metric.Duration)
}

func (a *MedianHandler) String() string {
	return fmt.Sprintf("[median] %s", metrics.Median(a.series))
}
