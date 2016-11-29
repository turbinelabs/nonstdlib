package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"fmt"
	"log"
	"runtime"
	"time"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

type DelayType string

const (
	ConstantDelayType    DelayType = "constant"
	ExponentialDelayType DelayType = "exponential"

	flagDefaultDelayType      = ExponentialDelayType
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
	// return the same Executor.
	Make(*log.Logger) Executor
}

// FromFlagsDefaults represents default values for Executor
// flags. Values are ignored if they are the zero value for their
// type.
type FromFlagsDefaults struct {
	DelayType      DelayType
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	MaxAttempts    int
	MaxQueueDepth  int
	Parallelism    int
	Timeout        time.Duration
	AttemptTimeout time.Duration
}

// Constructs a FromFlags with application-agnostic default flag
// values. Most callers should use NewFromFlagsWithDefaults.
func NewFromFlags(f *tbnflag.PrefixedFlagSet) FromFlags {
	return NewFromFlagsWithDefaults(f, FromFlagsDefaults{})
}

// Constructs a FromFlags with application-provided default flag
// values.
func NewFromFlagsWithDefaults(
	f *tbnflag.PrefixedFlagSet,
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
		fmt.Sprintf("retry delay type (%s or %s)", ExponentialDelayType, ConstantDelayType),
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

// If not overridden, the default delay type is exponential.
func (defaults FromFlagsDefaults) DefaultDelayType() DelayType {
	if defaults.DelayType != DelayType("") {
		return defaults.DelayType
	}
	return flagDefaultDelayType
}

// If not overridden, the default initial delay is 100 milliseconds.
func (defaults FromFlagsDefaults) DefaultInitialDelay() time.Duration {
	if defaults.InitialDelay != 0 {
		return defaults.InitialDelay
	}
	return flagDefaultInitialDelay
}

// If not overridden, the default maximum delay is 30 seconds.
func (defaults FromFlagsDefaults) DefaultMaxDelay() time.Duration {
	if defaults.MaxDelay != 0 {
		return defaults.MaxDelay
	}
	return flagDefaultMaxDelay
}

// If not overridden, the default max attempts is 8.
func (defaults FromFlagsDefaults) DefaultMaxAttempts() int {
	if defaults.MaxAttempts != 0 {
		return defaults.MaxAttempts
	}

	return flagDefaultMaxAttempts
}

// If not overridden, the default max queue depth is 20 times the
// number of system CPU cores.
func (defaults FromFlagsDefaults) DefaultMaxQueueDepth() int {
	if defaults.MaxQueueDepth != 0 {
		return defaults.MaxQueueDepth
	}

	return runtime.NumCPU() * 20
}

// If not overridden, the default parallelism is 2 times the number of
// system CPU cores.
func (defaults FromFlagsDefaults) DefaultParallelism() int {
	if defaults.Parallelism != 0 {
		return defaults.Parallelism
	}

	return runtime.NumCPU() * 2
}

// If not overridden, the default timeout is 0 (timeouts disabled).
func (defaults FromFlagsDefaults) DefaultTimeout() time.Duration {
	if defaults.Timeout != 0 {
		return defaults.Timeout
	}

	return flagDefaultTimeout
}

// If not overridden, the default attempt timeout is 0 (attempt
// timeouts disabled).
func (defaults FromFlagsDefaults) DefaultAttemptTimeout() time.Duration {
	if defaults.AttemptTimeout != 0 {
		return defaults.AttemptTimeout
	}

	return flagDefaultAttemptTimeout
}

type fromFlags struct {
	delayType      tbnflag.Choice
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

		ff.executor = NewRetryingExecutor(
			WithRetryDelayFunc(delayFunc),
			WithMaxAttempts(ff.maxAttempts),
			WithMaxQueueDepth(ff.maxQueueDepth),
			WithParallelism(ff.parallelism),
			WithTimeout(ff.timeout),
			WithAttemptTimeout(ff.attemptTimeout),
			WithLogger(log),
		)
	}

	return ff.executor
}
