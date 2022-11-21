package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"time"
)

type MaxHandler struct {
	result time.Duration
	record *parser.QueryRecord
}

func (a *MaxHandler) Init() MetricHandler {
	a.result = 0
	return a
}

func (a *MaxHandler) Process(metric *metrics.Metric) {
	if metric.Duration > a.result {
		a.result = metric.Duration
		a.record = metric.Record
	}
}

func (a *MaxHandler) String() string {
	if a.record != nil {
		return fmt.Sprintf("[max time] %s query[%s]", a.result, a.record)
	}
	return "[max] no valid data"
}
