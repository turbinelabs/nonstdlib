/*
Copyright 2018 Turbine Labs, Inc.

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

func testSignal(t *testing.T, f SignalFunc) (*procTestOutput, error) {
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

	return output, cmd.Wait()
}

func TestLoggingCmdKill(t *testing.T) {
	output, waitErr := testSignal(
		t,
		func(cmd *LoggingCmd) error {
			return cmd.Kill()
		},
	)

	assert.NonNil(t, waitErr)

	// expect no further output since kill cannot be trapped
	assert.False(t, output.checkCompleted())
	assert.HasSuffix(t, output.String(), "READY\n")
}

func TestLoggingCmdQuit(t *testing.T) {
	output, waitErr := testSignal(
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
