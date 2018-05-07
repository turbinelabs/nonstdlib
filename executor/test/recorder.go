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
	"sort"
	"sync/atomic"
	"time"
)

type resultType int

const (
	successResult resultType = iota
	failureResult
	timeoutResult
	tooManyCallsResult
	badResultType
	wrongJobResult
	unknownResult
)

func (r resultType) String() string {
	switch r {
	case successResult:
		return "success"
	case failureResult:
		return "failure"
	case timeoutResult:
		return "timeout"
	case tooManyCallsResult:
		return "too-many-calls"
	case badResultType:
		return "wrong-type"
	case wrongJobResult:
		return "wrong-job"
	case unknownResult:
		return "unknown"
	}

	panic(fmt.Sprintf("missing resultType (%d)", r))
}

type attempt struct {
	*job
	success bool
}

type callback struct {
	jobID  int64
	result resultType
}

type recorder struct {
	jobIDs    chan int64
	attempts  chan *attempt
	callbacks chan *callback

	quitChan chan struct{}
	doneChan chan struct{}
	updates  int64

	inflightJobs      map[int64]struct{}
	completedAttempts map[int64][]*attempt
	completedJobs     map[int64]*callback

	success       int64
	totalAttempts int64
	failures      int64
	timeouts      int64
	invalid       int64
}

func newRecorder() *recorder {
	return &recorder{
		jobIDs:            make(chan int64, 100),
		attempts:          make(chan *attempt, 100),
		callbacks:         make(chan *callback, 100),
		quitChan:          make(chan struct{}),
		doneChan:          make(chan struct{}),
		inflightJobs:      map[int64]struct{}{},
		completedAttempts: map[int64][]*attempt{},
		completedJobs:     map[int64]*callback{},
	}
}

func (r *recorder) start() {
	go r.run()
}

func (r *recorder) stop(d time.Duration) {
	timer := time.NewTimer(d)
	for {
		fmt.Println("waiting for quiesence")
		n := atomic.LoadInt64(&r.updates)
		<-timer.C
		afterN := atomic.LoadInt64(&r.updates)
		if n == afterN {
			break
		}
		timer.Reset(d)
	}

	close(r.quitChan)
	<-r.doneChan
}

func (r *recorder) run() {
	defer close(r.doneChan)

LOOP:
	for {
		var jobID int64
		select {
		case <-r.quitChan:
			break LOOP

		case callback := <-r.callbacks:
			jobID = callback.jobID
			if _, ok := r.completedJobs[jobID]; ok {
				fmt.Printf("ERROR: job completed more than once %d", jobID)
			}
			r.completedJobs[jobID] = callback

		case a := <-r.attempts:
			jobID = a.job.id
			if attempts, ok := r.completedAttempts[jobID]; ok {
				r.completedAttempts[jobID] = append(attempts, a)
			} else {
				attempts := make([]*attempt, 1, a.job.numFailures+1)
				attempts[0] = a
				r.completedAttempts[jobID] = attempts
			}

		case id := <-r.jobIDs:
			jobID = id
			if _, ok := r.inflightJobs[jobID]; ok {
				fmt.Printf("ERROR: restarted job %d", jobID)
			}
			r.inflightJobs[id] = struct{}{}
		}

		atomic.AddInt64(&r.updates, 1)

		r.checkJob(jobID)
	}

	for jobID := range r.inflightJobs {
		r.checkJob(jobID)
	}
	r.dump()
}

func (r *recorder) checkJob(jobID int64) {
	_, isInflight := r.inflightJobs[jobID]
	attempts, haveAttempts := r.completedAttempts[jobID]
	callback, hasCallback := r.completedJobs[jobID]

	dprintf(
		"\t\tjob %x: inflight? %t hasAttempts? %t, completed? %t\n",
		jobID,
		isInflight,
		haveAttempts,
		hasCallback,
	)

	if !isInflight || !haveAttempts || !hasCallback {
		return
	}

	//	fmt.Println(attempts)
	expectedAttempts := attempts[0].job.numFailures + 1
	if len(attempts) != expectedAttempts {
		dprintf(
			"\t\tjob %x: expected %d attempts, saw %d\n",
			jobID,
			expectedAttempts,
			len(attempts),
		)
		return
	}

	r.totalAttempts += int64(len(attempts))

	switch callback.result {
	case successResult:
		r.success++

	case failureResult:
		// TODO: should we have failed?
		r.failures++

	case timeoutResult:
		// TODO: should we have timed out
		r.timeouts++

	case tooManyCallsResult, badResultType, wrongJobResult, unknownResult:
		r.invalid++
		fmt.Printf("ERROR: %s\n", callback.result.String())
	}

	dprintf("\t\tjob %x: done\n", jobID)

	delete(r.inflightJobs, jobID)
	delete(r.completedAttempts, jobID)
	delete(r.completedJobs, jobID)
}

func (r *recorder) dump() {
	inflight := make([]int64, 0, len(r.inflightJobs))
	for jobID := range r.inflightJobs {
		inflight = append(inflight, jobID)
	}
	sort.Slice(inflight, func(i, j int) bool { return inflight[i] < inflight[j] })
	if len(inflight) > 0 {
		fmt.Println("\nInflight jobs:")
		for _, jobID := range inflight {
			fmt.Printf("  %d\n", jobID)
		}
	}

	attempted := []struct {
		jobID       int64
		numAttempts int
	}{}
	for jobID, attempts := range r.completedAttempts {
		attempted = append(
			attempted,
			struct {
				jobID       int64
				numAttempts int
			}{
				jobID:       jobID,
				numAttempts: len(attempts),
			},
		)
	}
	sort.Slice(
		attempted,
		func(i, j int) bool { return attempted[i].jobID < attempted[j].jobID },
	)
	if len(attempted) > 0 {
		fmt.Println("\nAttempted jobs:")
		for _, a := range attempted {
			fmt.Printf("  %d (%d times)", a.jobID, a.numAttempts)
		}
	}

	completed := make([]int64, 0, len(r.completedJobs))
	for jobID := range r.completedJobs {
		completed = append(completed, jobID)
	}
	sort.Slice(completed, func(i, j int) bool { return completed[i] < completed[j] })
	if len(completed) > 0 {
		fmt.Println("\nCompleted jobs:")
		for _, jobID := range completed {
			fmt.Printf("  %d", jobID)
		}
	}

	fmt.Println("\nResults:")
	fmt.Printf("  Successful runs: %d\n", r.success)
	fmt.Printf("  Total attempts: %d\n", r.totalAttempts)
	fmt.Printf("  Failed jobs: %d\n", r.failures)
	fmt.Printf("  Timed out jobs: %d\n", r.timeouts)
	fmt.Printf("  Invalid results: %d\n", r.invalid)
}
