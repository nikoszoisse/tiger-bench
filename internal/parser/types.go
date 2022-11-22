package parser

import (
	"fmt"
	"time"
)

// QueryRecord represents the record that we are going to build sql
type QueryRecord struct {
	Hostname  string
	StartTime time.Time
	EndTime   time.Time
	Line      int
}

// Result is the result of a csvReader after transforming row data
type Result struct {
	QueryRecord *QueryRecord
	Error       error
}

// String human-readable format
func (q *QueryRecord) String() string {
	return fmt.Sprintf(
		"host: %s, start_time: %s, end_time: %s (line %d)",
		q.Hostname, q.StartTime.Format(TimeLayout), q.EndTime.Format(TimeLayout), q.Line,
	)
}
