package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"time"
)

type MinHandler struct {
	result time.Duration
	record *parser.QueryRecord
}

func (a *MinHandler) Init() MetricHandler {
	const maxUint = ^uint(0)
	a.result = time.Duration(maxUint >> 1)

	return a
}

func (a *MinHandler) Process(metric *metrics.Metric) {
	if metric.Duration < a.result {
		a.result = metric.Duration
		a.record = metric.Record
	}
}

func (a *MinHandler) String() string {
	if a.record != nil {
		return fmt.Sprintf("[min time] %s query[%s]", a.result, a.record)
	}
	return "[min] no valid data"
}
