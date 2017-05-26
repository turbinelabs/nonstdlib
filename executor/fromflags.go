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

package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"log"
	"runtime"
	"time"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

// DelayType represents an algorithm for computing retry delays.
type DelayType string

const (
	// ConstantDelayType specifies a constant delay between
	// retries.
	ConstantDelayType DelayType = "constant"

	// ExponentialDelayType specifies an exponentially increasing
	// delay between retries.
	ExponentialDelayType DelayType = "exponential"

	flagDefaultDelayType      = ExponentialDelayType
	flagDefaultExperimental   = false
	flagDefaultInitialDelay   = 100 * time.Millisecond
	flagDefaultMaxDelay       = 30 * time.Second
	flagDefaultMaxAttempts    = 8
	flagDefaultTimeout        = 0 * time.Second
	flagDefaultAttemptTimeout = 0 * time.Second
)

// FromFlags validates and constructs an Executor from command line
// flags.
type FromFlags interface {
	// Returns the configured Executor. Multiple invocations
	// return the same Executor even if the arguments
	// change. DiagnosticsCallback may be nil.
	Make(*log.Logger) Executor
}

// FromFlagsDefaults represents default values for Executor
// flags. Values are ignored if they are the zero value for their
// type.
type FromFlagsDefaults struct {
	Experimental   *bool
	DelayType      DelayType
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	MaxAttempts    int
	MaxQueueDepth  int
	Parallelism    int
	Timeout        time.Duration
	AttemptTimeout time.Duration
}

// NewFromFlags constructs a FromFlags with application-agnostic
// default flag values. Most callers should use NewFromFlagsWithDefaults.
func NewFromFlags(f tbnflag.FlagSet) FromFlags {
	return NewFromFlagsWithDefaults(f, FromFlagsDefaults{})
}

// NewFromFlagsWithDefaults constructs a FromFlags with
// application-provided default flag values.
func NewFromFlagsWithDefaults(
	f tbnflag.FlagSet,
	defaults FromFlagsDefaults,
) FromFlags {
	delayTypeChoice :=
		tbnflag.NewChoice(string(ConstantDelayType), string(ExponentialDelayType)).
			WithDefault(string(defaults.DefaultDelayType()))

	ff := &fromFlags{
		delayType: delayTypeChoice,
	}

	f.Var(
		&ff.delayType,
		"delay-type",
		"Specifies the retry delay type.",
	)

	f.BoolVar(
		&ff.experimental,
		"experimental",
		defaults.DefaultExperimental(),
		"Enables an experiment goroutine-based executor.",
	)

	f.DurationVar(
		&ff.initialDelay,
		"delay",
		defaults.DefaultInitialDelay(),
		"Specifies the initial delay for the exponential delay type. "+
			"Specifies the delay for constant delay type.",
	)

	f.DurationVar(
		&ff.maxDelay,
		"max-delay",
		defaults.DefaultMaxDelay(),
		"Specifies the maximum delay for the exponential delay type. "+
			"Ignored for the constant delay type.",
	)

	f.IntVar(
		&ff.maxAttempts,
		"max-attempts",
		defaults.DefaultMaxAttempts(),
		"Specifies the maximum number of attempts made, inclusive of the original attempt.",
	)

	f.IntVar(
		&ff.maxQueueDepth,
		"max-queue",
		defaults.DefaultMaxQueueDepth(),
		"Specifies the maximum number of attempts that may be queued before new "+
			"attempts are blocked.",
	)

	f.IntVar(
		&ff.parallelism,
		"parallelism",
		defaults.DefaultParallelism(),
		"Specifies the maximum number of concurrent attempts running.",
	)

	f.DurationVar(
		&ff.timeout,
		"timeout",
		defaults.DefaultTimeout(),
		"Specifies the default timeout for actions. A timeout of 0 means no timeout.",
	)

	f.DurationVar(
		&ff.attemptTimeout,
		"attempt-timeout",
		defaults.DefaultAttemptTimeout(),
		"Specifies the default timeout for individual action attempts. A timeout of 0 "+
			"means no timeout.",
	)

	return ff
}

// DefaultDelayType returns the default delay type. If not overridden
// the default delay type is ExponentialDelayType.
func (defaults FromFlagsDefaults) DefaultDelayType() DelayType {
	if defaults.DelayType != DelayType("") {
		return defaults.DelayType
	}
	return flagDefaultDelayType
}

// DefaultExperimental returns the experimental setting. If not
// overridden the default delay type is false.
func (defaults FromFlagsDefaults) DefaultExperimental() bool {
	if defaults.Experimental != nil {
		return *defaults.Experimental
	}
	return flagDefaultExperimental
}

// DefaultInitialDelay returns the default initial delay. If not
// overridden, the default initial delay is 100 milliseconds.
func (defaults FromFlagsDefaults) DefaultInitialDelay() time.Duration {
	if defaults.InitialDelay != 0 {
		return defaults.InitialDelay
	}
	return flagDefaultInitialDelay
}

// DefaultMaxDelay returns the default maximum delay. If not
// overridden, the default maximum delay is 30 seconds.
func (defaults FromFlagsDefaults) DefaultMaxDelay() time.Duration {
	if defaults.MaxDelay != 0 {
		return defaults.MaxDelay
	}
	return flagDefaultMaxDelay
}

// DefaultMaxAttempts returns the default maximum number of
// attempts. If not overridden, the default max attempts is 8.
func (defaults FromFlagsDefaults) DefaultMaxAttempts() int {
	if defaults.MaxAttempts != 0 {
		return defaults.MaxAttempts
	}

	return flagDefaultMaxAttempts
}

// DefaultMaxQueueDepth returns the default maximum queue depth. If
// not overridden, the default max queue depth is 20 times the number
// of system CPU cores.
func (defaults FromFlagsDefaults) DefaultMaxQueueDepth() int {
	if defaults.MaxQueueDepth != 0 {
		return defaults.MaxQueueDepth
	}

	return runtime.NumCPU() * 20
}

// DefaultParallelism returns the default parallelism. If not
// overridden, the default parallelism is 2 times the number of system
// CPU cores.
func (defaults FromFlagsDefaults) DefaultParallelism() int {
	if defaults.Parallelism != 0 {
		return defaults.Parallelism
	}

	return runtime.NumCPU() * 2
}

// DefaultTimeout returns the default global timeout. If not
// overridden, the default timeout is 0 (timeouts disabled).
func (defaults FromFlagsDefaults) DefaultTimeout() time.Duration {
	if defaults.Timeout != 0 {
		return defaults.Timeout
	}

	return flagDefaultTimeout
}

// DefaultAttemptTimeout returns the default per-attempt timeout. If
// not overridden, the default attempt timeout is 0 (attempt timeouts
// disabled).
func (defaults FromFlagsDefaults) DefaultAttemptTimeout() time.Duration {
	if defaults.AttemptTimeout != 0 {
		return defaults.AttemptTimeout
	}

	return flagDefaultAttemptTimeout
}

type fromFlags struct {
	delayType      tbnflag.Choice
	experimental   bool
	initialDelay   time.Duration
	maxDelay       time.Duration
	maxAttempts    int
	maxQueueDepth  int
	parallelism    int
	timeout        time.Duration
	attemptTimeout time.Duration

	executor Executor
}

func (ff *fromFlags) Make(log *log.Logger) Executor {
	if ff.executor == nil {
		var delayFunc DelayFunc
		switch DelayType(ff.delayType.String()) {
		case ExponentialDelayType:
			delayFunc = NewExponentialDelayFunc(ff.initialDelay, ff.maxDelay)
		case ConstantDelayType:
			delayFunc = NewConstantDelayFunc(ff.initialDelay)
		}

		options := []Option{
			WithRetryDelayFunc(delayFunc),
			WithMaxAttempts(ff.maxAttempts),
			WithMaxQueueDepth(ff.maxQueueDepth),
			WithParallelism(ff.parallelism),
			WithTimeout(ff.timeout),
			WithAttemptTimeout(ff.attemptTimeout),
			WithLogger(log),
		}

		if ff.experimental {
			ff.executor = NewGoroutineExecutor(options...)
		} else {
			ff.executor = NewRetryingExecutor(options...)
		}
	}

	return ff.executor
}
