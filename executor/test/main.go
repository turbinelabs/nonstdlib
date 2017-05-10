package test

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/turbinelabs/nonstdlib/executor"
	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

const (
	defaultMinDelay    = 50 * time.Millisecond
	defaultMaxDelay    = 250 * time.Millisecond
	defaultMinFailures = 0
	defaultMaxFailures = 3
)

var (
	debug = false
)

func dprintf(f string, args ...interface{}) {
	if debug {
		fmt.Printf(f, args...)
	}
}

func dprintln(args ...interface{}) {
	if debug {
		fmt.Println(args...)
	}
}

// BakeTestCLI configures and parses command line flags and runs a
// bake test of the nonstdlib/executor.
func BakeTestCLI() {
	ff, err := newFromFlags()
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if err := ff.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := ff.Make()
	exec := config.exec
	debug = config.debug
	defer exec.Stop()

	recorder := newRecorder()
	recorder.start()

	generators := make([]*generator, config.numJobGenerators)
	for i := range generators {
		generators[i] = &generator{
			exec:        exec,
			id:          int32(i),
			rate:        config.rate / float64(config.numJobGenerators),
			minDelay:    config.minDelay,
			maxDelay:    config.maxDelay,
			minFailures: config.minFailures,
			maxFailures: config.maxFailures,
		}

		generators[i].init()
		generators[i].start(recorder)
	}

	tickerStep := 1 * time.Second
	if config.stopAfter > 10*time.Second {
		tickerStep = 5 * time.Second
	}
	if config.stopAfter > 2*time.Minute {
		tickerStep = 30 * time.Second
	}
	if config.stopAfter > 10*time.Minute {
		tickerStep = 5 * time.Minute
	}
	if config.stopAfter > 1*time.Hour {
		tickerStep = 15 * time.Minute
	}
	ticker := time.NewTicker(tickerStep)
	go func() {
		left := config.stopAfter
		for range ticker.C {
			left -= tickerStep
			fmt.Printf("(%s remaining)\n", left.String())
		}
	}()
	time.Sleep(config.stopAfter)
	ticker.Stop()

	fmt.Println("Stopping generator(s)...")
	for _, g := range generators {
		g.stop()
	}

	fmt.Println("Stopping recorder...")
	delay := config.timeout
	if delay == 0 {
		delay = config.attemptTimeout * time.Duration(config.maxAttempts)
		if delay == 0 {
			delay = config.maxDelay * time.Duration(config.maxFailures+1)
		}
	}
	recorder.stop(delay)

	fmt.Println("Stopping executor...")
	exec.Stop()
}

type bakeTestConfig struct {
	debug bool

	numJobGenerators int
	rate             float64
	stopAfter        time.Duration
	minDelay         time.Duration
	maxDelay         time.Duration
	minFailures      int
	maxFailures      int

	exec           executor.Executor
	delayType      string
	attemptTimeout time.Duration
	timeout        time.Duration
	maxAttempts    int
}

type fromFlags struct {
	flagSet       *flag.FlagSet
	execFromFlags executor.FromFlags

	config bakeTestConfig
}

func newFromFlags() (*fromFlags, error) {
	ff := &fromFlags{
		flagSet: flag.NewFlagSet("bake-test", flag.ContinueOnError),
	}
	tbnFlagSet := tbnflag.Wrap(ff.flagSet)

	tbnFlagSet.BoolVar(
		&ff.config.debug,
		"debug",
		false,
		"Enable debug output",
	)

	tbnFlagSet.IntVar(
		&ff.config.numJobGenerators,
		"num-gen",
		4,
		"Number of job generating goroutines.",
	)
	tbnFlagSet.Float64Var(
		&ff.config.rate,
		"rate",
		1.0,
		"Number of jobs/second",
	)
	tbnFlagSet.DurationVar(
		&ff.config.stopAfter,
		"stop-after",
		time.Minute,
		"Sets the test run time.",
	)
	tbnFlagSet.DurationVar(
		&ff.config.minDelay,
		"min-delay",
		defaultMinDelay,
		"Minimum attempt delay.",
	)
	tbnFlagSet.DurationVar(
		&ff.config.maxDelay,
		"max-delay",
		defaultMaxDelay,
		"Maximum attempt delay.",
	)
	tbnFlagSet.IntVar(
		&ff.config.minFailures,
		"min-failures",
		defaultMinFailures,
		"Minimum attempt failures (per job).",
	)
	tbnFlagSet.IntVar(
		&ff.config.maxFailures,
		"max-failures",
		defaultMaxFailures,
		"Maximum attempt failures (per job).",
	)

	ff.execFromFlags = executor.NewFromFlags(tbnFlagSet.Scope("exec", "Executor"))

	if err := ff.flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return ff, nil
}

func (ff *fromFlags) Validate() error {
	durationFlags := map[string]time.Duration{}
	intFlags := map[string]int{}
	delayType := ""

	ff.flagSet.VisitAll(func(f *flag.Flag) {
		if !strings.HasPrefix(f.Name, "exec.") {
			return
		}
		if g, isGetter := f.Value.(flag.Getter); isGetter {
			switch v := g.Get().(type) {
			case time.Duration:
				durationFlags[f.Name] = v
			case int:
				intFlags[f.Name] = v
			default:
				if f.Name == "exec.delay-type" {
					delayType = f.Value.String()
				}
			}
		} else {
			panic(fmt.Sprintf("cannot get value for %s", f.Name))
		}
	})

	ff.config.attemptTimeout = durationFlags["exec.attempt-timeout"]
	ff.config.timeout = durationFlags["exec.timeout"]
	ff.config.maxAttempts = intFlags["exec.max-attempts"]

	if ff.config.maxDelay < ff.config.minDelay {
		ff.config.minDelay, ff.config.maxDelay = ff.config.maxDelay, ff.config.minDelay
	}

	if ff.config.attemptTimeout > 0 && ff.config.minDelay > ff.config.attemptTimeout {
		fmt.Printf(
			"WARNING: all jobs will timeout (min delay %s > attempt timeout %s)\n",
			ff.config.minDelay,
			ff.config.attemptTimeout,
		)
	}

	if ff.config.timeout > 0 && ff.config.minDelay > ff.config.timeout {
		fmt.Printf(
			"WARNING: all jobs will timeout (min delay %s > global timeout %s)\n",
			ff.config.minDelay,
			ff.config.timeout,
		)
	}

	if ff.config.maxFailures < ff.config.minFailures {
		ff.config.minFailures, ff.config.maxFailures = ff.config.maxFailures, ff.config.minFailures
	}

	if ff.config.minFailures > ff.config.maxAttempts {
		fmt.Printf(
			"WARNING: all jobs will fail (min failures %d > max attempts %d)\n",
			ff.config.minFailures,
			ff.config.maxAttempts,
		)
	}

	return nil
}

func (ff *fromFlags) Make() *bakeTestConfig {
	ff.config.exec = ff.execFromFlags.Make(nil)
	return &ff.config
}
