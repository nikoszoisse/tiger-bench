package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/nikoszoisse/tiger-bench/pkg/logger"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

const TimeLayout = "2006-01-02 15:04:05"

const csvColumns = 3
const hostNamePos = 0
const startTimePos = 1
const endTimePos = 2

// ReadCsvFile will try to open the file then spawn a goroutine, and return the output channel where  results get published
// output channel is buffered based on free memory space
func ReadCsvFile(filePath string) chan *Result {
	if len(filePath) > 0 && !strings.HasSuffix(filePath, ".csv") {
		log.Fatal("Not valid input file, only '.csv' files are allowed." + filePath)
	}

	f, err := openFile(filePath)
	if f == nil || err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	pageSize, err := calcPageSize(f)
	if err != nil {
		logger.DefaultLogger.Error("Could not determine pageSize, read line by line. error %s", err.Error())
		pageSize = 0 // unbuffered
	}
	// Create the channel
	outputChannel := make(chan *Result, pageSize)

	// Spawn the routine
	go readCsvFileLazy(f, outputChannel)

	return outputChannel
}

// calcPageSize calc how many times the file rows we are reading can be loaded into "free" memory (fitCount)
// this way we have mem based back pressure while reading.
// We divide the size of a record from the "free memory" resulting the fitCount
// If the fit-count is > lines/2 then we return a pageSize = lines
// If the fit-count is < lines/2 or 0 then we return pageSize = fitCount
func calcPageSize(f *os.File) (int, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	lines, err := countLines(f)
	queryRecordSize := uint(unsafe.Sizeof(QueryRecord{}))

	fitCount := (m.Sys - m.TotalAlloc) / uint64(queryRecordSize)
	if fitCount > uint64(lines/2) {
		return int(lines), err
	}

	return int(fitCount), err
}

func countLines(f *os.File) (uint, error) {
	defer func() {
		// rewind to the beginning of the file
		f.Seek(0, io.SeekStart)
	}()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return uint(count), nil

		case err != nil:
			return uint(count), err
		}
	}
}

// openFile if filePath is not empty otherwise open stdin
func openFile(filePath string) (*os.File, error) {
	// Choose data source
	if len(filePath) == 0 {
		return os.Stdin, nil
	}
	return os.Open(filePath)
}

// readCsvFileLazy read csv from the given filePath, or reading from the stdin when filePath is empty
// reads a row and mapping it to a QueryRecord and then publish it as Result to the resultChannel
// In case of invalid rows a Result with Result.Error is returned
func readCsvFileLazy(file *os.File, resultChannel chan *Result) {
	defer func() {
		close(resultChannel)
		file.Close()
	}()
	csvReader := csv.NewReader(file)

	var timeReadBegin *time.Time
	//Start Reading the rows by line -> async publish it
	for {
		// Read the row
		row, err := csvReader.Read()
		if timeReadBegin == nil {
			t := time.Now()
			timeReadBegin = &t
		}
		// Check if EOF then stop reading
		if errors.Is(err, io.EOF) {
			return
		} else if len(row) != csvColumns { // Check if the row has the right column schema
			resultChannel <- &Result{nil, errors.New(fmt.Sprintf("Invalid row format data: %s", row))}
		} else {
			// Take the line of query for debug purposes
			line, _ := csvReader.FieldPos(0)

			// Build the QueryRecord from row
			queryRecord, err := buildRecord(row, line, *timeReadBegin)
			if err != nil {
				resultChannel <- &Result{nil, err}
				continue
			}

			//Publish Result
			resultChannel <- &Result{queryRecord, nil}
		}
	}
}

// buildRecord will transform row data to a QueryRecord
func buildRecord(record []string, line int, begin time.Time) (*QueryRecord, error) {
	startTime, err := time.Parse(TimeLayout, record[startTimePos])
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(TimeLayout, record[endTimePos])
	if err != nil {
		return nil, err
	}
	return &QueryRecord{
		Hostname:    record[hostNamePos],
		StartTime:   startTime,
		EndTime:     endTime,
		Line:        line,
		CreatedTime: begin.Add(time.Now().Sub(begin)),
	}, nil
}
