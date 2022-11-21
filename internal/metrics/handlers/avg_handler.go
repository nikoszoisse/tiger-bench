package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"time"
)

type AvgHandler struct {
	counter int64
	sum     time.Duration
	result  time.Duration
}

func (a *AvgHandler) String() string {
	return fmt.Sprintf("[avg] %s", a.result)
}

func (a *AvgHandler) Process(metric *metrics.Metric) {
	a.sum += metric.Duration
	a.counter++
	a.result = time.Duration(a.sum.Nanoseconds() / a.counter)
}

func (a *AvgHandler) Init() MetricHandler {
	a.counter = 0
	a.sum = 0
	a.result = 0

	return a
}
