package proc

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"syscall"
)

// A wrapper around os.exec.Cmd that logs stdout and stderr to a
// log.Logger.
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

// Send a quit signal to the process. Signaling an unstarted process
// does not produce an error.
func (c *LoggingCmd) Quit() error {
	if c.Process != nil {
		return c.Process.Signal(syscall.SIGQUIT)
	}

	return nil
}

// Retrieve the process id, if started. Otherwise returns -1.
func (c *LoggingCmd) Pid() int {
	if c.Process != nil {
		return c.Process.Pid
	}

	return -1
}

// Constructs a LoggingCmd that logs to the default logger. The
// cmd and args parameters are passed directly to
// os.exec.Command. See that method for details.
func DefaultLoggingCommand(cmd string, args ...string) *LoggingCmd {
	return LoggingCommand(nil, cmd, args...)
}

// Constructs a LoggingCmd that logs to the given logger. The cmd
// and args parameters are passed directly to os.exec.Command. See
// that method for details.
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
