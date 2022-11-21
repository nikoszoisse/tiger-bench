package worker

import (
	"context"
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/config"
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"github.com/nikoszoisse/tiger-bench/pkg/logger"
	"github.com/serialx/hashring"
	"math"
	"strings"
)

type WorkersRouter struct {
	workers   map[string]*worker
	hashRing  *hashring.HashRing
	ctx       context.Context
	processCh chan *parser.QueryRecord
}

func NewWorkersRouter(ctx context.Context, count int, dbConfig *config.DbConfig, workerResultHandler WorkerOutputHandler) *WorkersRouter {
	router := &WorkersRouter{
		workers:   make(map[string]*worker),
		ctx:       ctx,
		processCh: make(chan *parser.QueryRecord, 2*count),
	}

	router.hashRing = hashring.New([]string{})
	for i := 0; i < count; i++ {
		worker := spawnWorker(ctx, dbConfig, fmt.Sprintf("worker-%d", i), int64(count), workerResultHandler)

		router.workers[worker.name] = worker
		router.hashRing = router.hashRing.AddNode(worker.name)
	}

	keys := make([]string, 0, len(router.workers))
	multiplier := math.Log(float64(len(router.workers)))
	for k := range router.workers {
		for i := 0; i <= int(multiplier*float64(len(router.workers))); i++ {
			keys = append(keys, fmt.Sprintf("%d:%s", i, k))
		}
	}
	router.hashRing = hashring.New(keys)

	router.run()

	return router
}

func (r *WorkersRouter) Process(record *parser.QueryRecord) {
	r.processCh <- record
}

func (r *WorkersRouter) run() {
	go func() {
		defer func() {
			close(r.processCh)
		}()
		for {
			select {
			case <-r.ctx.Done():
				logger.DefaultLogger.Debug("WorkersRouter shutting down.")
				return
			case record := <-r.processCh:
				if selectedWorker, ok := r.hashRing.GetNode(record.Hostname); ok {
					actualWorkerName := strings.Split(selectedWorker, ":")[1]
					r.workers[actualWorkerName].inputChannel <- record
				} else {
					logger.DefaultLogger.Error("Could not route record: %s", record)
				}
			}
		}
	}()
}
