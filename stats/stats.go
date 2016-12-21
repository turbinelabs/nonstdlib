/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
