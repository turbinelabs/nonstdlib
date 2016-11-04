// Provides simple console logging to Stderr, configurable through a FlagSet.
// There are three levels of logging:
//  - none
//  - error
//  - debug
// The default log level is error
package console

import (
	"flag"
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
// error or debug, therwise it returns a no-op Logger.
func Error() *log.Logger {
	if *logLevelChoice.Choice == noneLevel {
		return nullLogger
	}
	return errorLogger
}

// Debug returns a Logger to Stderr prefixed with "[debug]" if the log level is
// debug, therwise it returns a no-op Logger.
func Debug() *log.Logger {
	if *logLevelChoice.Choice == debugLevel {
		return debugLogger
	}
	return nullLogger
}

// Init binds the log level to a flag in the given FlagSet
func Init(fs *flag.FlagSet) {
	fs.Var(
		&logLevelChoice,
		"console.level",
		"if none, log nothing. if error, log only errors. if debug, include debug logs.",
	)
}
