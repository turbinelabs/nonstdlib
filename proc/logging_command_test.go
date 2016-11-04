package proc

import (
	"strings"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

type SignalFunc func(*LoggingCmd) error

func TestLoggingCmdStdout(t *testing.T) {
	logger, output := makeProcTestOutput()
	cmd := LoggingCommand(logger, "echo", "xypdq")

	err := cmd.Start()
	assert.Nil(t, err)

	err = cmd.Wait()
	assert.Nil(t, err)

	assert.Equal(t, output.String(), "[echo stdout] xypdq\n")
}

func TestLoggingCmdStderr(t *testing.T) {
	logger, output := makeProcTestOutput()
	cmd := LoggingCommand(logger, "sh", "-c", "echo xypdq >/dev/null 1>&2")

	err := cmd.Start()
	assert.Nil(t, err)

	err = cmd.Wait()
	assert.Nil(t, err)

	assert.Equal(t, output.String(), "[sh stderr] xypdq\n")
}

func testSignal(t *testing.T, f SignalFunc) (error, *procTestOutput) {
	logger, output := makeProcTestOutput()
	exec, args := procTestAndArgs(10)
	cmd := LoggingCommand(logger, exec, args...)

	err := f(cmd)
	assert.Nil(t, err)

	err = cmd.Start()
	assert.Nil(t, err)

	for !output.checkReady() {
		time.Sleep(10 * time.Millisecond)
	}

	err = f(cmd)
	assert.Nil(t, err)

	return cmd.Wait(), output
}

func TestLoggingCmdKill(t *testing.T) {
	waitErr, output := testSignal(
		t,
		func(cmd *LoggingCmd) error {
			return cmd.Kill()
		},
	)

	assert.NonNil(t, waitErr)

	// expect no further output since kill cannot be trapped
	assert.False(t, output.checkCompleted())
	assert.True(t, strings.HasSuffix(output.String(), "READY\n"))
}

func TestLoggingCmdQuit(t *testing.T) {
	waitErr, output := testSignal(
		t,
		func(cmd *LoggingCmd) error {
			return cmd.Quit()
		},
	)

	assert.Nil(t, waitErr)
	assert.True(t, output.checkSignaled())
}

func TestLoggingCmdPid(t *testing.T) {
	logger, _ := makeProcTestOutput()
	exec, args := procTestAndArgs(10)
	cmd := LoggingCommand(logger, exec, args...)

	pid := cmd.Pid()
	assert.Equal(t, pid, -1)

	err := cmd.Start()
	assert.Nil(t, err)

	pid = cmd.Pid()
	if pid <= 0 {
		t.Errorf("invalid pid %d", pid)
	}

	err = cmd.Kill()
	assert.Nil(t, err)

	err = cmd.Wait()
	assert.NonNil(t, err)
}
