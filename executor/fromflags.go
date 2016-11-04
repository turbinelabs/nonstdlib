package executor

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"fmt"
	"log"
	"runtime"
	"time"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

const (
	constantDelayType    = "constant"
	exponentialDelayType = "exponential"
	defaultDelayType     = exponentialDelayType
)

// FromFlags validates and constructs an Executor from command line
// flags.
type FromFlags interface {
	// Returns the configured Executor. Multiple invocations
	// return the same Executor.
	Make(*log.Logger) Executor
}

func NewFromFlags(f *tbnflag.PrefixedFlagSet) FromFlags {
	ff := &fromFlags{
		delayType: tbnflag.NewChoice(constantDelayType, exponentialDelayType).
			WithDefault(exponentialDelayType),
	}

	f.Var(
		&ff.delayType,
		"delay-type",
		fmt.Sprintf("retry delay type (%s or %s)", exponentialDelayType, constantDelayType),
	)

	f.DurationVar(
		&ff.initialDelay,
		"delay",
		100*time.Millisecond,
		"Specifies the initial delay for the exponential delay type. "+
			"Specifies the delay for constant delay type.",
	)

	f.DurationVar(
		&ff.maxDelay,
		"max-delay",
		30*time.Second,
		"Specifies the maximum delay for the exponential delay type. "+
			"Ignored for the constant delay type.",
	)

	f.IntVar(
		&ff.maxAttempts,
		"max-attempts",
		8,
		"Specifies the maximum number of attempts made, inclusive of the original attempt.",
	)

	f.IntVar(
		&ff.maxQueueDepth,
		"max-queue",
		runtime.NumCPU()*20,
		"Specifies the maximum number of attempts that may be queued before new "+
			"attempts are blocked.",
	)

	f.IntVar(
		&ff.parallelism,
		"parallelism",
		runtime.NumCPU()*2,
		"Specifies the maximum number of concurrent attempts running.",
	)

	f.DurationVar(
		&ff.timeout,
		"timeout",
		0,
		"Specifies the default timeout for actions. A timeout of 0 means no timeout.",
	)
	return ff
}

type fromFlags struct {
	delayType     tbnflag.Choice
	initialDelay  time.Duration
	maxDelay      time.Duration
	maxAttempts   int
	maxQueueDepth int
	parallelism   int
	timeout       time.Duration

	executor Executor
}

func (ff *fromFlags) Make(log *log.Logger) Executor {
	if ff.executor == nil {
		var delayFunc DelayFunc
		switch ff.delayType.String() {
		case exponentialDelayType:
			delayFunc = NewExponentialDelayFunc(ff.initialDelay, ff.maxDelay)
		case constantDelayType:
			delayFunc = NewConstantDelayFunc(ff.initialDelay)
		}

		ff.executor = NewRetryingExecutor(
			WithRetryDelayFunc(delayFunc),
			WithMaxAttempts(ff.maxAttempts),
			WithMaxQueueDepth(ff.maxQueueDepth),
			WithParallelism(ff.parallelism),
			WithTimeout(ff.timeout),
			WithLogger(log),
		)
	}

	return ff.executor
}
