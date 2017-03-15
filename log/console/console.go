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

// Package console provides simple console logging to Stderr, configurable
// through a FlagSet.
// There are three levels of logging:
//  - none
//  - error
//  - debug
// The default log level is error.
//
// Executable using the console package should include exactly one call to
// Init() with the flag.FlagSet used to configure the executable, passed
// prior to the FlagSet being parsed. Subsequently, calls to Error() and
// Debug() will produce output to os.Stderr, based on the log-level configured.
package console

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

const (
	noneLevel  = "none"
	errorLevel = "error"
	infoLevel  = "info"
	debugLevel = "debug"

	defaultLevel = infoLevel
)

const (
	noneOrdinal int = iota
	errorOrdinal
	infoOrdinal
	debugOrdinal

	defaultOrdinal = infoOrdinal
)

var (
	logFlags = log.LstdFlags | log.LUTC

	errorLogger = log.New(os.Stderr, "[error] ", logFlags)
	infoLogger  = log.New(os.Stderr, "[info] ", logFlags)
	debugLogger = log.New(os.Stderr, "[debug] ", logFlags)
	nullLogger  = log.New(ioutil.Discard, "", 0)

	logLevelChoice = tbnflag.NewChoice(
		debugLevel,
		infoLevel,
		errorLevel,
		noneLevel,
	).WithDefault(defaultLevel)

	logLevelOrder = map[string]int{
		noneLevel:  noneOrdinal,
		errorLevel: errorOrdinal,
		infoLevel:  infoOrdinal,
		debugLevel: debugOrdinal,
	}
)

func logLevel() int {
	choice := logLevelChoice.Choice
	if choice == nil {
		return defaultOrdinal
	}

	if level, ok := logLevelOrder[*choice]; ok {
		return level
	}

	return defaultOrdinal
}

// Error returns a Logger to Stderr prefixed with "[error]" if the log level is
// error, info, or debug, otherwise it returns a no-op Logger.
func Error() *log.Logger {
	if logLevel() < errorOrdinal {
		return nullLogger
	}
	return errorLogger
}

// Info returns a Logger to Stderr prefixed with "[info]" if the log
// level is info or error, otherwise it returns a no-op Logger.
func Info() *log.Logger {
	if logLevel() < infoOrdinal {
		return nullLogger
	}
	return infoLogger
}

// Debug returns a Logger to Stderr prefixed with "[debug]" if the log level is
// debug, otherwise it returns a no-op Logger.
func Debug() *log.Logger {
	if logLevel() < debugOrdinal {
		return nullLogger
	}

	return debugLogger
}

// Init binds the log level to a flag in the given FlagSet.
func Init(fs tbnflag.FlagSet) {
	fs.Var(
		&logLevelChoice,
		"console.level",
		"Selects the log `level` for console logs messages.",
	)
}

// ToWriteCloser wraps a *log.Logger and provides an implementation of
// io.WriteCloser
func ToWriteCloser(l *log.Logger) io.WriteCloser {
	reader, writer := io.Pipe()
	buffered := bufio.NewReader(reader)
	done := make(chan bool)
	go func() {
		for {
			line, err := buffered.ReadString('\n')
			// eg EOF
			if err != nil {
				done <- true
				return
			}
			l.Output(1, line)
		}
	}()
	return &waitWriteCloser{writer, done}
}

type waitWriteCloser struct {
	io.WriteCloser
	done chan bool
}

func (c *waitWriteCloser) Close() error {
	err := c.WriteCloser.Close()
	// wait until the log writing function has signaled done before returning
	<-c.done
	return err
}
