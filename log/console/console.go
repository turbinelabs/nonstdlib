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
	"io/ioutil"
	"log"
	"os"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
)

const (
	errorLevel = "error"
	debugLevel = "debug"
	noneLevel  = "none"
)

var (
	errorLogger    = log.New(os.Stderr, "[error] ", log.LstdFlags)
	debugLogger    = log.New(os.Stderr, "[debug] ", log.LstdFlags)
	nullLogger     = log.New(ioutil.Discard, "", 0)
	logLevelChoice = tbnflag.NewChoice(debugLevel, errorLevel, noneLevel).WithDefault(errorLevel)
)

// Error returns a Logger to Stderr prefixed with "[error]" if the log level is
// error or debug, otherwise it returns a no-op Logger.
func Error() *log.Logger {
	if *logLevelChoice.Choice == noneLevel {
		return nullLogger
	}
	return errorLogger
}

// Debug returns a Logger to Stderr prefixed with "[debug]" if the log level is
// debug, otherwise it returns a no-op Logger.
func Debug() *log.Logger {
	if *logLevelChoice.Choice == debugLevel {
		return debugLogger
	}
	return nullLogger
}

// Init binds the log level to a flag in the given FlagSet.
func Init(fs tbnflag.FlagSet) {
	fs.Var(
		&logLevelChoice,
		"console.level",
		"if none, log nothing. if error, log only errors. if debug, include debug logs.",
	)
}
