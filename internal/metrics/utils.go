package metrics

import (
	"context"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"time"
)

type Number interface {
	constraints.Float | constraints.Integer
}

// WithTrack will track the time when a method take place
func WithTrack(ctx context.Context, method func(ctx context.Context) error) (time.Duration, error) {
	start := time.Now()
	err := method(ctx)
	return time.Since(start), err
}

// Median implementation
func Median[T Number](series []T) T {
	seriesSnap := make([]T, len(series))
	copy(seriesSnap, series)

	slices.Sort(seriesSnap)

	var median T
	l := len(seriesSnap)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (seriesSnap[l/2-1] + seriesSnap[l/2]) / 2.0
	} else {
		median = seriesSnap[l/2]
	}

	return median
}
