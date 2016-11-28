package stats

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import "time"

// Stats represents a simple interface for forwarding arbitrary stats.
type Stats interface {
	Inc(string, int64) error
	Gauge(string, int64) error
	TimingDuration(string, time.Duration) error
}

// Creates a do-nothing Stats.
func NewNoopStats() Stats {
	return &noopStats{}
}

type noopStats struct{}

func (_ *noopStats) Inc(_ string, _ int64) error                    { return nil }
func (_ *noopStats) Gauge(_ string, _ int64) error                  { return nil }
func (_ *noopStats) TimingDuration(_ string, _ time.Duration) error { return nil }

// Creates a Stats implementation that forwards all calls to an
// underlying Stats using goroutines. All methods always return nil.
func NewAsyncStats(s Stats) Stats {
	return &asyncStats{underlying: s}
}

type asyncStats struct {
	underlying Stats
}

func (a *asyncStats) Inc(name string, value int64) error {
	go a.underlying.Inc(name, value)
	return nil
}

func (a *asyncStats) Gauge(name string, value int64) error {
	go a.underlying.Gauge(name, value)
	return nil
}

func (a *asyncStats) TimingDuration(name string, value time.Duration) error {
	go a.underlying.TimingDuration(name, value)
	return nil
}
