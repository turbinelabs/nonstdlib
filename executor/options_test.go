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

package executor

import (
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func TestWithLogger(t *testing.T) {
	log := log.NewNoopLogger()

	exec := &commonExec{}
	WithLogger(log)(exec)
	assert.SameInstance(t, exec.log, log)
}

func TestWithDiagnostics(t *testing.T) {
	diag := NewNoopDiagnosticsCallback()

	exec := &commonExec{}
	WithDiagnostics(diag)(exec)
	assert.SameInstance(t, exec.diag, diag)
}

func TestWithRetryDelayFunc(t *testing.T) {
	d := NewConstantDelayFunc(0)

	exec := &commonExec{}
	WithRetryDelayFunc(d)(exec)
	assert.SameInstance(t, exec.delay, d)
}

func TestWithMaxAttempts(t *testing.T) {
	exec := &commonExec{}

	WithMaxAttempts(0)(exec)
	assert.Equal(t, exec.maxAttempts, 1)

	WithMaxAttempts(100)(exec)
	assert.Equal(t, exec.maxAttempts, 100)

	WithMaxAttempts(-100)(exec)
	assert.Equal(t, exec.maxAttempts, 1)
}

func TestWithParallelism(t *testing.T) {
	exec := &commonExec{}

	WithParallelism(0)(exec)
	assert.Equal(t, exec.parallelism, 1)

	WithParallelism(100)(exec)
	assert.Equal(t, exec.parallelism, 100)

	WithParallelism(-100)(exec)
	assert.Equal(t, exec.parallelism, 1)
}

func TestWithTimeout(t *testing.T) {
	exec := &commonExec{}

	WithTimeout(0 * time.Second)(exec)
	assert.Equal(t, exec.timeout, 0*time.Second)

	WithTimeout(100 * time.Millisecond)(exec)
	assert.Equal(t, exec.timeout, 100*time.Millisecond)

	WithTimeout(-100 * time.Millisecond)(exec)
	assert.Equal(t, exec.timeout, 0*time.Second)
}
