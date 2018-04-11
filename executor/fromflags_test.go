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
	"runtime"
	"testing"
	"time"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func TestFromFlags(t *testing.T) {
	log := log.NewNoopLogger()

	flagSet := tbnflag.NewTestFlagSet()

	ff := NewFromFlags(flagSet.Scope("exec", "whatever"))
	assert.NonNil(t, ff)

	ffImpl := ff.(*fromFlags)

	expectedParallelism := runtime.NumCPU() * 2

	assert.Equal(t, ffImpl.delayType.String(), string(ExponentialDelayType))
	assert.Equal(t, ffImpl.initialDelay, 100*time.Millisecond)
	assert.Equal(t, ffImpl.maxDelay, 30*time.Second)
	assert.Equal(t, ffImpl.maxAttempts, 8)
	assert.Equal(t, ffImpl.parallelism, expectedParallelism)
	assert.Nil(t, ffImpl.executor)

	diag := NewNoopDiagnosticsCallback()

	exec := ff.Make(log)
	exec.SetDiagnosticsCallback(diag)
	assert.SameInstance(t, exec, ffImpl.executor)
	assert.SameInstance(t, ff.Make(log), exec)

	commonImpl, ok := exec.(*commonExec)
	assert.True(t, ok)

	expExecImpl, ok := commonImpl.impl.(*goroutineExecImpl)
	assert.True(t, ok)

	assert.NonNil(t, expExecImpl.sem)
	assert.Equal(t, cap(expExecImpl.sem), expectedParallelism)
	assert.Equal(t, commonImpl.parallelism, expectedParallelism)
	assert.Equal(t, commonImpl.maxAttempts, 8)
	assert.NonNil(t, commonImpl.delay)
	assert.Equal(t, commonImpl.delay(1), 100*time.Millisecond)
	assert.Equal(t, commonImpl.delay(100000), 30*time.Second)
	assert.Equal(t, commonImpl.timeout, 0*time.Second)
	assert.Equal(t, commonImpl.attemptTimeout, 0*time.Second)
	assert.SameInstance(t, commonImpl.log, log)
	assert.SameInstance(t, commonImpl.diag, diag)

	exec.Stop()

	ffImpl.executor = nil

	flagSet.Parse([]string{
		"-exec.delay-type=constant",
		"-exec.delay=1s",
		"-exec.max-delay=5s",
		"-exec.max-attempts=4",
		"-exec.parallelism=99",
		"-exec.timeout=100ms",
		"-exec.attempt-timeout=10ms",
	})

	assert.Equal(t, ffImpl.delayType.String(), string(ConstantDelayType))
	assert.Equal(t, ffImpl.initialDelay, time.Second)
	assert.Equal(t, ffImpl.maxDelay, 5*time.Second)
	assert.Equal(t, ffImpl.maxAttempts, 4)
	assert.Equal(t, ffImpl.parallelism, 99)
	assert.Equal(t, ffImpl.timeout, 100*time.Millisecond)
	assert.Equal(t, ffImpl.attemptTimeout, 10*time.Millisecond)

	expectedParallelism = 99

	exec = ff.Make(nil)
	assert.SameInstance(t, exec, ffImpl.executor)

	commonImpl, ok = exec.(*commonExec)
	assert.True(t, ok)

	expExecImpl, ok = commonImpl.impl.(*goroutineExecImpl)
	assert.True(t, ok)

	assert.NonNil(t, expExecImpl.sem)
	assert.Equal(t, cap(expExecImpl.sem), expectedParallelism)
	assert.Equal(t, commonImpl.parallelism, expectedParallelism)
	assert.Equal(t, commonImpl.maxAttempts, 4)
	assert.NonNil(t, commonImpl.delay)
	assert.Equal(t, commonImpl.delay(1), 1*time.Second)
	assert.Equal(t, commonImpl.delay(100000), 1*time.Second)
	assert.Equal(t, commonImpl.timeout, 100*time.Millisecond)
	assert.Nil(t, commonImpl.log)
	_, ok = commonImpl.diag.(*noopDiagnosticsCallback)
	assert.True(t, ok)

	exec.Stop()
}

func TestFromFlagsWithDefaults(t *testing.T) {
	prefixedFlagSet := tbnflag.NewTestFlagSet().Scope("exec", "whatever")
	ff := NewFromFlagsWithDefaults(prefixedFlagSet, FromFlagsDefaults{})
	ffImpl := ff.(*fromFlags)

	assert.Equal(t, ffImpl.delayType.String(), string(ExponentialDelayType))
	assert.Equal(t, ffImpl.initialDelay, flagDefaultInitialDelay)
	assert.Equal(t, ffImpl.maxDelay, flagDefaultMaxDelay)
	assert.Equal(t, ffImpl.maxAttempts, flagDefaultMaxAttempts)
	assert.Equal(t, ffImpl.parallelism, 2*runtime.NumCPU())
	assert.Equal(t, ffImpl.timeout, 0*time.Second)
	assert.Equal(t, ffImpl.attemptTimeout, 0*time.Second)

	prefixedFlagSet = tbnflag.NewTestFlagSet().Scope("exec", "whatever")
	ff = NewFromFlagsWithDefaults(
		prefixedFlagSet,
		FromFlagsDefaults{
			DelayType:      ConstantDelayType,
			InitialDelay:   1 * time.Second,
			MaxDelay:       2 * time.Second,
			MaxAttempts:    3,
			Parallelism:    5,
			Timeout:        6 * time.Second,
			AttemptTimeout: 7 * time.Millisecond,
		},
	)
	ffImpl = ff.(*fromFlags)

	assert.Equal(t, ffImpl.delayType.String(), string(ConstantDelayType))
	assert.Equal(t, ffImpl.initialDelay, 1*time.Second)
	assert.Equal(t, ffImpl.maxDelay, 2*time.Second)
	assert.Equal(t, ffImpl.maxAttempts, 3)
	assert.Equal(t, ffImpl.parallelism, 5)
	assert.Equal(t, ffImpl.timeout, 6*time.Second)
	assert.Equal(t, ffImpl.attemptTimeout, 7*time.Millisecond)
}
