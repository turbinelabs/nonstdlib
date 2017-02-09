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

package proc

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"syscall"
)

// LoggingCmd is a wrapper around os.exec.Cmd that logs stdout and
// stderr to a log.Logger.
type LoggingCmd struct {
	*exec.Cmd
}

// Kill the process. Killing an unstarted process does not produce an
// error.
func (c *LoggingCmd) Kill() error {
	if c.Process != nil {
		return c.Process.Kill()
	}

	return nil
}

// Quit sends a quit signal to the process. Signaling an unstarted process
// does not produce an error.
func (c *LoggingCmd) Quit() error {
	if c.Process != nil {
		return c.Process.Signal(syscall.SIGQUIT)
	}

	return nil
}

// Pid retrieves the process id, if started. Otherwise returns -1.
func (c *LoggingCmd) Pid() int {
	if c.Process != nil {
		return c.Process.Pid
	}

	return -1
}

// DefaultLoggingCommand constructs a LoggingCmd that logs to the
// default logger. The cmd and args parameters are passed directly to
// os.exec.Command. See that method for details.
func DefaultLoggingCommand(cmd string, args ...string) *LoggingCmd {
	return LoggingCommand(nil, cmd, args...)
}

// LoggingCommand constructs a LoggingCmd that logs to the given
// logger. The cmd and args parameters are passed directly to
// os.exec.Command. See that method for details.
func LoggingCommand(logger *log.Logger, cmd string, args ...string) *LoggingCmd {
	name := path.Base(cmd)

	underlying := exec.Command(cmd, args...)

	var stdout, stderr LogWriter
	if logger == nil {
		stdout = NewDefaultLogWriter(fmt.Sprintf("[%s stdout] ", name))
		stderr = NewDefaultLogWriter(fmt.Sprintf("[%s stderr] ", name))
	} else {
		stdout = NewLogWriter(logger, fmt.Sprintf("[%s stdout] ", name))
		stderr = NewLogWriter(logger, fmt.Sprintf("[%s stderr] ", name))
	}
	underlying.Stdout = &stdout
	underlying.Stderr = &stderr

	return &LoggingCmd{underlying}
}
