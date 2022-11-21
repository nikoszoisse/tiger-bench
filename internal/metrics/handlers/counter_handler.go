package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
)

type CounterHandler struct {
	counter int64
}

func (a *CounterHandler) String() string {
	return fmt.Sprintf("[total queries #] %d", a.counter)
}

func (a *CounterHandler) Process(metric *metrics.Metric) {
	a.counter++
}

func (a *CounterHandler) Init() MetricHandler {
	a.counter = 0

	return a
}
