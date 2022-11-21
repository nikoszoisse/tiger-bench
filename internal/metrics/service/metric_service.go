package service

import (
	"context"
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"github.com/nikoszoisse/tiger-bench/internal/metrics/handlers"
	"github.com/nikoszoisse/tiger-bench/pkg/logger"
)

type MetricService struct {
	inputChannel    chan *metrics.Metric
	ctx             context.Context
	cancelFn        context.CancelFunc
	metricHandlers  []handlers.MetricHandler
	processCallback func(metric *metrics.Metric)
}

// NewMetricsService register handlers that for every metric that is coming is calculated
// the capacity of the metrics channel is recommended to be the equals as the Reader's output buffer channel
func NewMetricsService(ctx context.Context, capacity int, afterProcess func(metric *metrics.Metric)) *MetricService {
	metricsContext, cancelFn := context.WithCancel(ctx)
	service := &MetricService{
		inputChannel: make(chan *metrics.Metric, capacity),
		ctx:          metricsContext,
		cancelFn:     cancelFn,
		metricHandlers: []handlers.MetricHandler{
			(&handlers.CounterHandler{}).Init(),
			(&handlers.AvgHandler{}).Init(),
			(&handlers.MedianHandler{}).Init(),
			(&handlers.StdHandler{}).Init(),
			(&handlers.MinHandler{}).Init(),
			(&handlers.MaxHandler{}).Init(),
		},
		processCallback: afterProcess,
	}

	service.run()

	return service
}

// Process TODO add ctx here
func (m *MetricService) Process(result *metrics.Metric) {
	m.inputChannel <- result
}

func (m *MetricService) run() {
	go func() {
		defer func() {
			close(m.inputChannel)
			m.cancelFn()
		}()

		for {
			select {
			case <-m.ctx.Done():
				logger.DefaultLogger.Debug("MetricService shutting down.")
				return
			case result := <-m.inputChannel:
				for _, handler := range m.metricHandlers {
					handler.Process(result)
				}
				m.processCallback(result)
			}
		}
	}()
}

func (m *MetricService) Results() string {
	result := "\n=====Stats=====\n"
	for _, handler := range m.metricHandlers {
		result += fmt.Sprintf("%s \n", handler)
	}
	result += "\n===============\n"
	return result
}
