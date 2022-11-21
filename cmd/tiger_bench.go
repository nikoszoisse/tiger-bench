package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/nikoszoisse/tiger-bench/internal/config"
	"github.com/nikoszoisse/tiger-bench/internal/metrics"
	"github.com/nikoszoisse/tiger-bench/internal/metrics/service"
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"github.com/nikoszoisse/tiger-bench/internal/worker"
	"github.com/nikoszoisse/tiger-bench/pkg/logger"
	"log"
	"runtime"
	"sync"
	"time"
)

var appConfig *config.AppConfig

func Run() {
	appConfig = readConfigs()

	// Start Reading CSV File and consume content through channel
	queryRecordChannel := parser.ReadCsvFile(appConfig.FilePath)

	metricsService, totalDuration := benchmark(queryRecordChannel)

	fmt.Println(metricsService.Results())
	fmt.Printf("[total time] %s\n\n", totalDuration)
	fmt.Println("===============")
}

func readConfigs() *config.AppConfig {
	// Parse cli arguments
	verbose := flag.Bool("v", false, "Print verbose logs")
	filePath := flag.String("file", "", "query params file path")
	dbUrl := flag.String("db", "", "database url e.g. localhost:5432,your_db,user,password")
	numOfWorkers := flag.Int("workers", runtime.NumCPU(), "num of worker, default: Num Of CPUs")

	flag.Parse()

	dbConfig := config.NewDbConfig(*dbUrl)
	if dbConfig == nil {
		log.Fatal("[ERROR] database url is invalid or missing: ", *dbUrl)
	}

	if *verbose {
		logger.DefaultLogger.SetLevel(logger.DebugLevel)
	}

	return &config.AppConfig{
		Verbose:      *verbose,
		FilePath:     *filePath,
		DbConfig:     dbConfig,
		WorkersCount: *numOfWorkers,
	}
}

func benchmark(queryRecordChannel chan *parser.Result) (*service.MetricService, time.Duration) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Set up a wait-group in order to wait for workers and metrics to finish processing
	waitGroup := sync.WaitGroup{}

	// When the Metric service process a metric then we're Done from the wait-group
	metricsService := service.NewMetricsService(ctx, cap(queryRecordChannel), func(metric *metrics.Metric) {
		// metric processed
		waitGroup.Done()
	})

	// Spawn a worker router of given numOfWorkers
	// When the record is processed successfully we send its metrics data to metric service
	// If record is dirty then we .Done in wait-group
	router := worker.NewWorkersRouter(ctx, appConfig.WorkersCount, appConfig.DbConfig,
		func(workerId worker.WorkerID, record *parser.QueryRecord, duration time.Duration) {
			if record != nil {
				logger.DefaultLogger.Debug(workerId+": "+"Record (%s, %s, %s) took %s \n", record.Hostname, record.StartTime, record.EndTime, duration)
				metricsService.Process(&metrics.Metric{Record: record, Duration: duration})
			} else {
				waitGroup.Done()
			}
		})

	duration, _ := metrics.WithTrack(context.Background(), func(ctx context.Context) error {
		// Process the QueryRecordsChannel
		// Route the record
		for record := range queryRecordChannel {
			if record.QueryRecord == nil || record.Error != nil {
				logger.DefaultLogger.Error("Record %s err %s \n", record.QueryRecord, record.Error)
			} else {
				waitGroup.Add(1)
				router.Process(record.QueryRecord)
			}
		}

		waitGroup.Wait()

		return nil
	})

	return metricsService, duration
}
