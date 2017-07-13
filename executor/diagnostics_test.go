package executor

import (
	"io"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/turbinelabs/nonstdlib/arrays/indexof"
	tbntime "github.com/turbinelabs/nonstdlib/time"
	"github.com/turbinelabs/test/assert"
	"github.com/turbinelabs/test/log"
)

func TestAttemptResultString(t *testing.T) {
	for r := AttemptResult(0); r < attemptUnknown; r++ {
		assert.NotEqual(t, r.String(), "")
		assert.NotEqual(t, r.String(), "AttemptUnknown")
	}

	assert.Equal(t, attemptUnknown.String(), "AttemptUnknown")
	assert.Equal(t, AttemptResult(-1).String(), "AttemptUnknown")
	assert.Equal(t, AttemptResult(numAttemptResults+1).String(), "AttemptUnknown")
}

func TestAttemptResultValid(t *testing.T) {
	for r := AttemptResult(0); r < attemptUnknown; r++ {
		assert.True(t, r.Valid())
	}

	assert.False(t, attemptUnknown.Valid())
	assert.False(t, AttemptResult(-1).Valid())
	assert.False(t, AttemptResult(numAttemptResults+1).Valid())
}

func TestForEachAttemptResult(t *testing.T) {
	values := []int{}
	names := []string{}
	ForEachAttemptResult(
		func(r AttemptResult) {
			values = append(values, int(r))
			names = append(names, r.String())
		},
	)

	assert.Equal(t, len(values), numAttemptResults)
	assert.True(t, sort.IntsAreSorted(values))

	sort.Strings(names)
	for i := 1; i < len(names); i++ {
		assert.NotEqual(t, names[i-1], names[i])
	}

	assert.Equal(t, indexof.String(names, attemptUnknown.String()), indexof.NotFound)
}

func TestCountedDuration(t *testing.T) {
	cd := &countedDuration{}
	msg, ok := cd.format("prefix")
	assert.False(t, ok)
	assert.Equal(t, msg, "prefix: 0")

	cd.add(time.Minute)
	assert.Equal(t, cd.count, int64(1))
	assert.Equal(t, time.Duration(cd.totalDuration), time.Minute)
	assert.Equal(t, time.Duration(cd.maxDuration), time.Minute)

	cd.add(30 * time.Second)
	cd.add(30 * time.Second)
	assert.Equal(t, cd.count, int64(3))
	assert.Equal(t, time.Duration(cd.totalDuration), 2*time.Minute)
	assert.Equal(t, time.Duration(cd.maxDuration), time.Minute)

	msg, ok = cd.format("prefix")
	assert.True(t, ok)
	assert.Equal(t, msg, "prefix: 3 (avg 40s; max 1m0s)")
}

func TestCountedDurationMultiThreaded(t *testing.T) {
	cd := &countedDuration{}

	start := make(chan struct{})
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(3)
	go func() {
		defer waitGroup.Done()

		select {
		case <-start:
		}

		for i := 0; i < 1000; i++ {
			cd.add(time.Second)
		}
		for i := 0; i < 1000; i++ {
			cd.add(time.Minute)
		}
	}()

	go func() {
		defer waitGroup.Done()

		select {
		case <-start:
		}

		for i := 0; i < 2000; i++ {
			cd.add(time.Duration(i*50) * time.Millisecond)
		}
	}()

	go func() {
		defer waitGroup.Done()

		select {
		case <-start:
		}

		for i := 0; i < 2000; i++ {
			cd.add(time.Duration(i*50) * time.Millisecond)
		}
	}()

	close(start)
	waitGroup.Wait()

	assert.Equal(t, cd.count, int64(6000))
	assert.Equal(
		t,
		time.Duration(cd.totalDuration),
		(1000*time.Second)+(1000*time.Minute)+2*(1000*1999*50*time.Millisecond),
	)
	assert.Equal(t, time.Duration(cd.maxDuration), 1999*50*time.Millisecond)
}

func TestCountedDurationsByResult(t *testing.T) {
	cdbr := newCountedDurationsByResult()
	s, any := cdbr.format("prefix ")
	assert.Equal(t, len(s), 0)
	assert.False(t, any)

	cdbr.add(AttemptSuccess, time.Second)
	cdbr.add(AttemptSuccess, 2*time.Second)
	cdbr.add(AttemptError, time.Second)
	cdbr.add(AttemptError, time.Second)

	s, any = cdbr.format("prefix ")
	assert.ArrayEqual(
		t,
		s,
		[]string{
			"prefix AttemptSuccess: 2 (avg 1.5s; max 2s)",
			"prefix AttemptError: 2 (avg 1s; max 1s)",
		},
	)
	assert.True(t, any)

	cdbr.add(AttemptTimeout, time.Second)
	cdbr.add(AttemptGlobalTimeout, time.Second)
	cdbr.add(AttemptCancellation, time.Second)
	cdbr.add(attemptUnknown, time.Second)

	s, any = cdbr.format("prefix2 ")
	assert.ArrayEqual(
		t,
		s,
		[]string{
			"prefix2 AttemptSuccess: 2 (avg 1.5s; max 2s)",
			"prefix2 AttemptTimeout: 1 (avg 1s; max 1s)",
			"prefix2 AttemptGlobalTimeout: 1 (avg 1s; max 1s)",
			"prefix2 AttemptCancellation: 1 (avg 1s; max 1s)",
			"prefix2 AttemptError: 2 (avg 1s; max 1s)",
			"prefix2 AttemptUnknown: 1 (avg 1s; max 1s)",
		},
	)
	assert.True(t, any)
}

func TestNewLoggingDiagnosticsCallback(t *testing.T) {
	logger := log.NewNoopLogger()
	ldc := NewLoggingDiagnosticsCallback(logger, time.Minute)

	ldcImpl := ldc.(*loggingDiagnosticsCallback)
	assert.SameInstance(t, ldcImpl.logger, logger)
	assert.Equal(t, ldcImpl.period, time.Minute)
	assert.NonNil(t, ldcImpl.time)
	assert.NonNil(t, ldcImpl.quit)
	assert.NonNil(t, ldcImpl.data)
	assert.Nil(t, ldcImpl.Close())
}

func TestLoggingDiagnosticsCallback(t *testing.T) {
	logger := log.NewNoopLogger()
	ldc := NewLoggingDiagnosticsCallback(logger, time.Hour)
	defer ldc.(io.Closer).Close()

	ldc.TaskStarted(5)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptCompleted(AttemptTimeout, 1*time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptResult(100), 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 2*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptResult(200), 1*time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)

	ldcImpl := ldc.(*loggingDiagnosticsCallback)
	assert.Equal(t, ldcImpl.data.tasksStarted, int64(5))
	assert.Equal(t, ldcImpl.data.tasksCompleted[AttemptSuccess].count, int64(4))
	assert.Equal(t, ldcImpl.data.tasksCompleted[attemptUnknown].count, int64(1))
	assert.Equal(t, ldcImpl.data.attemptsStarted.count, int64(6))
	assert.Equal(t, ldcImpl.data.attemptsCompleted[AttemptSuccess].count, int64(4))
	assert.Equal(t, ldcImpl.data.attemptsCompleted[AttemptTimeout].count, int64(1))
	assert.Equal(t, ldcImpl.data.attemptsCompleted[attemptUnknown].count, int64(1))
	assert.Equal(t, ldcImpl.data.callbacks.count, int64(2))
}

func TestLoggingDiagnosticsCallbackLog(t *testing.T) {
	logger, buffer := log.NewBufferLogger()
	ldc := NewLoggingDiagnosticsCallback(logger, time.Hour)
	defer ldc.(io.Closer).Close()

	ldc.(*loggingDiagnosticsCallback).log()

	expected := `Executor Diagnostics
tasks started: 0
`
	assert.Equal(t, buffer.String(), expected)
	buffer.Reset()

	ldc.TaskStarted(1)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptStarted(1 * time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.AttemptCompleted(AttemptError, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptSuccess, 1*time.Millisecond)
	ldc.TaskCompleted(AttemptTimeout, 1*time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)
	ldc.CallbackDuration(1 * time.Millisecond)

	ldc.(*loggingDiagnosticsCallback).log()

	expected = `Executor Diagnostics
tasks started: 1
tasks completed, AttemptSuccess: 4 (avg 1ms; max 1ms)
tasks completed, AttemptTimeout: 1 (avg 1ms; max 1ms)
attempts started: 2 (avg 1ms; max 1ms)
attempts completed, AttemptSuccess: 3 (avg 1ms; max 1ms)
attempts completed, AttemptError: 1 (avg 1ms; max 1ms)
callbacks: 5 (avg 1ms; max 1ms)
`
	assert.Equal(t, buffer.String(), expected)
}

func TestLoggingDiagnosticsCallbackLogPeriodically(t *testing.T) {
	tbntime.WithCurrentTimeFrozen(func(cs tbntime.ControlledSource) {
		logger, ch := log.NewChannelLogger(100)
		ldc := newLoggingDiagnosticsCallback(logger, time.Minute, cs)
		defer ldc.(io.Closer).Close()

		for cs.TriggerAllTimers() == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		assert.Equal(t, <-ch, "Executor Diagnostics\n")
		assert.Equal(t, <-ch, "tasks started: 0\n")

		ldc.TaskStarted(1)

		for cs.TriggerAllTimers() == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		assert.Equal(t, <-ch, "Executor Diagnostics\n")
		assert.Equal(t, <-ch, "tasks started: 1\n")

		ldc.TaskStarted(2)
		ldc.TaskStarted(1)

		for cs.TriggerAllTimers() == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		assert.Equal(t, <-ch, "Executor Diagnostics\n")
		assert.Equal(t, <-ch, "tasks started: 3\n")
	})
}
