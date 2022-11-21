package handlers

import (
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"math"
	"time"
)

type StdHandler struct {
	counter     int64
	sum         float64
	sumSquared  float64
	result      time.Duration
	resultError time.Duration
}

func (a *StdHandler) String() string {
	return fmt.Sprintf("[std] %s  with error: %s", a.result, a.resultError)
}

func (a *StdHandler) Process(metric *metrics.Metric) {
	a.sum += metric.Duration.Seconds()
	a.sumSquared += metric.Duration.Seconds() * metric.Duration.Seconds()
	a.counter++
	avg := a.sum / float64(a.counter)
	avgSq := a.sumSquared / float64(a.counter) // calculate the mean/std.dev
	a.result = time.Duration(float64(time.Second) * math.Sqrt(avgSq-avg*avg))
	a.resultError = time.Duration(float64(a.result.Nanoseconds()) / math.Sqrt(float64(a.counter)))
}

func (a *StdHandler) Init() MetricHandler {
	a.counter = 0
	a.sum = 0
	a.result = 0

	return a
}
