package worker

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nikoszoisse/tiger-bench/internal/config"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"github.com/nikoszoisse/tiger-bench/pkg/logger"
	"time"
)

type WorkerID = string
type WorkerOutputHandler = func(workerId WorkerID, record *parser.QueryRecord, duration time.Duration)

type worker struct {
	name          string
	inputChannel  chan *parser.QueryRecord
	dbConfig      *config.DbConfig
	db            *sql.Conn
	ctx           context.Context
	outputHandler WorkerOutputHandler
}

func spawnWorker(ctx context.Context, dbConfig *config.DbConfig, name string, bufferSize int64, handler WorkerOutputHandler) *worker {
	worker := &worker{
		name:          name,
		inputChannel:  make(chan *parser.QueryRecord, bufferSize),
		ctx:           ctx,
		outputHandler: handler,
		dbConfig:      dbConfig,
	}

	worker.run()
	return worker
}

func (w *worker) run() {
	dbCon, err := tryToRetrieveDBConnection(w.dbConfig)
	if err != nil {
		panic(fmt.Sprintf("worker %s, failed to connect to the dataSource: %s", w.name, w.dbConfig))
	}
	w.db = dbCon

	go func() {
		defer func() {
			close(w.inputChannel)
			w.db.Close()
		}()
		for {
			select {
			case <-w.ctx.Done():
				logger.DefaultLogger.Debug("%s shutting down.", w.name)
				return
			case record := <-w.inputChannel:
				if record != nil {
					logger.DefaultLogger.Debug("[worker %s] processing %+v \n", w.name, record)
					// build & execute the query
					//keep duration since query creation
					duration, err := metrics.WithTrack(w.ctx, func(ctx context.Context) error {
						return w.queryRecord(record)
					})

					if err != nil {
						w.outputHandler(w.name, nil, duration)
					} else {
						// TODO consider passing ct instead of name
						w.outputHandler(w.name, record, duration)
					}
				}
			}
		}
	}()
}

func (w *worker) queryRecord(record *parser.QueryRecord) error {
	// generate
	generatedSql := generateSql(record)
	// query
	rows, err := w.db.QueryContext(w.ctx, generatedSql)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	return err
}

func generateSql(record *parser.QueryRecord) string {
	return fmt.Sprintf(`SELECT time_bucket('1 minutes', ts) AS one_min, max(usage), min(usage)
		FROM cpu_usage
		WHERE host='%s' and ts between '%s' and '%s'
		GROUP BY one_min
		ORDER BY one_min DESC;`,
		record.Hostname, record.StartTime.Format(parser.TimeLayout), record.EndTime.Format(parser.TimeLayout))
}

func tryToRetrieveDBConnection(dbConfig *config.DbConfig) (*sql.Conn, error) {
	db, err := sql.Open("postgres", dbConfig.String())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db.Conn(context.Background())
}
