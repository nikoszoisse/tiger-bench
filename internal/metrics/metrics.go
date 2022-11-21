package metrics

import (
	"github.com/nikoszoisse/tiger-bench/internal/parser"
	"time"
)

type Metric struct {
	Record   *parser.QueryRecord
	Duration time.Duration
}
