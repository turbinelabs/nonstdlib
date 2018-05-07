/*
Copyright 2018 Turbine Labs, Inc.

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

package test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/turbinelabs/nonstdlib/executor"
	tbnmath "github.com/turbinelabs/nonstdlib/math"
)

type generator struct {
	exec executor.Executor

	id          int32
	rate        float64
	minDelay    time.Duration
	maxDelay    time.Duration
	minFailures int
	maxFailures int

	quitChan chan struct{}
	doneChan chan struct{}
}

func (g *generator) init() {
	g.quitChan = make(chan struct{})
	g.doneChan = make(chan struct{})
}

func (g *generator) start(r *recorder) {
	go g.run(r)
}

func (g *generator) stop() {
	close(g.quitChan)
	<-g.doneChan
}

func (g *generator) mkDelays(rng *rand.Rand, n int) []time.Duration {
	delays := make([]time.Duration, 0, n)

	delayRange := int64(g.maxDelay - g.minDelay)
	for ; n > 0; n-- {
		delayNanos := int64(g.minDelay)
		if delayRange > 0 {
			delayNanos += rng.Int63n(int64(delayRange + 1))
		}
		delays = append(delays, time.Duration(delayNanos))
	}
	return delays
}

func (g *generator) run(r *recorder) {
	defer close(g.doneChan)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	delayGen, err := tbnmath.NewPoissonDistributionWithRand(g.rate, rng)
	if err != nil {
		fmt.Printf("generator error: %s\n", err.Error())
		return
	}
	timer := time.NewTimer(time.Hour)

	jobID := int32(0)

	failureRange := g.maxFailures - g.minFailures

	for {
		delay := delayGen.Next()
		timer.Reset(delay)

		select {
		case <-g.quitChan:
			return
		case <-timer.C:
			// do a thing
		}

		numFailures := g.minFailures
		if failureRange != 0 {
			numFailures += rng.Intn(failureRange + 1)
		}
		job := &job{
			id:          (int64(g.id) << 32) | int64(jobID),
			numFailures: numFailures,
			delays:      g.mkDelays(rng, numFailures+1),
			recorder:    r,
		}
		jobID++

		r.jobIDs <- job.id

		g.exec.Exec(job.Go, job.Callback)
	}
}
