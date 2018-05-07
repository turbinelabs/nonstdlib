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
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/turbinelabs/nonstdlib/executor"
)

var (
	errJobFailed               = errors.New("failure")
	errJobExecutedTooManyTimes = errors.New("too many executions")
)

type job struct {
	id          int64
	exec        int
	numFailures int
	delays      []time.Duration
	recorder    *recorder
}

func (j *job) String() string {
	return fmt.Sprintf("job %x: %d (%v)", j.id, j.numFailures, j.delays)
}

func (j *job) Go(ctxt context.Context) (interface{}, error) {
	dprintf("\tjob %x: running (%d of %d)\n", j.id, j.exec+1, j.numFailures+1)
	attemptResult := &attempt{job: j}
	defer func() {
		j.recorder.attempts <- attemptResult
	}()

	if j.exec >= len(j.delays) {
		dprintf("\tjob %x: too many\n", j.id)
		return nil, errJobExecutedTooManyTimes
	}

	delay := j.delays[j.exec]
	j.exec++
	if delay > 0 {
		dprintf("\tjob %x: delay %s\n", j.id, delay.String())
		timer := time.NewTimer(delay)
		select {
		case <-ctxt.Done():
			return nil, ctxt.Err()
		case <-timer.C:
			// continue
		}
	}

	if j.exec > j.numFailures {
		dprintf("\tjob %x: succeed\n", j.id)
		attemptResult.success = true
		return j.id, nil
	}

	dprintf("\tjob %x: fail\n", j.id)
	return nil, errJobFailed
}

func (j *job) Callback(try executor.Try) {
	dprintf("\tjob: %x callback\n", j.id)
	result := unknownResult
	if try.IsError() {
		err := try.Error()
		if err == errJobFailed {
			result = failureResult
		} else if strings.Contains(err.Error(), "timeout") {
			result = timeoutResult
		}
	} else {
		if jobID, ok := try.Get().(int64); ok {
			if jobID == j.id {
				result = successResult
			} else {
				result = wrongJobResult
			}
		} else {
			result = badResultType
		}
	}

	j.recorder.callbacks <- &callback{jobID: j.id, result: result}
}
