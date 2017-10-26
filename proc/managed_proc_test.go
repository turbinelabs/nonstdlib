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
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

var (
	testExe, testArgs = procTestAndArgs(60)

	expectedCmdArgs = append([]string{testExe}, testArgs...)
)

type managedProcTest struct {
	proc   ManagedProc
	output *procTestOutput
	errors chan error
	wg     sync.WaitGroup
}

func makeManagedProcTest(t *testing.T) *managedProcTest {
	test := &managedProcTest{
		errors: make(chan error, 1),
	}

	test.wg.Add(1)
	onExit := func(err error) {
		test.errors <- err
		test.wg.Done()
	}

	logger, procTestOutput := makeProcTestOutput()

	proc := NewManagedProc(testExe, testArgs, logger, onExit)
	err := proc.Start()
	assert.Nil(t, err)
	assert.NonNil(t, proc)
	assert.True(t, proc.Running())

	for !procTestOutput.checkReady() {
		time.Sleep(10 * time.Millisecond)
	}

	procArgs := proc.(*managedProc).Args
	assert.ArrayEqual(t, procArgs, expectedCmdArgs)

	test.proc = proc
	test.output = procTestOutput

	return test
}

func expectExitSignal(t *testing.T, test *managedProcTest, wantSignal syscall.Signal) {
	var procErr error
	select {
	case e := <-test.errors:
		procErr = e
	default:
		procErr = nil
	}

	if !assert.NonNil(t, procErr) {
		return
	}

	exitErr, ok := procErr.(*exec.ExitError)
	if !ok {
		assert.Tracing(t).Errorf("process did not exit with an exec.ExitError")
	}

	waitStatus, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		assert.Tracing(t).Errorf("system dependent exit status is not a syscall.WaitStatus")
	}

	if waitStatus.Signaled() {
		assert.Equal(t, waitStatus.Signal(), wantSignal)
	} else {
		assert.Tracing(t).Errorf(
			"expected signal %s (0x%x), but process was not signaled: wait status %+v",
			wantSignal.String(),
			wantSignal,
			waitStatus)
	}
}

func TestManagedProc(t *testing.T) {
	test := makeManagedProcTest(t)

	assert.Nil(t, test.proc.Kill())

	test.wg.Wait()
	assert.False(t, test.proc.Running())
	expectExitSignal(t, test, syscall.SIGKILL)
}

func TestManagedProcHangup(t *testing.T) {
	test := makeManagedProcTest(t)

	assert.Nil(t, test.proc.Hangup())

	test.wg.Wait()
	assert.False(t, test.proc.Running())
	assert.True(t, test.output.checkSignaled())
}

func TestManagedProcTerm(t *testing.T) {
	test := makeManagedProcTest(t)

	assert.Nil(t, test.proc.Term())

	test.wg.Wait()
	assert.False(t, test.proc.Running())
	assert.True(t, test.output.checkSignaled())
}

func TestManagedProcUsr1(t *testing.T) {
	test := makeManagedProcTest(t)

	assert.Nil(t, test.proc.Usr1())

	test.wg.Wait()
	assert.False(t, test.proc.Running())
	assert.True(t, test.output.checkSignaled())
}

func TestManagedProcQuit(t *testing.T) {
	test := makeManagedProcTest(t)

	assert.Nil(t, test.proc.Quit())

	test.wg.Wait()
	assert.False(t, test.proc.Running())
	assert.True(t, test.output.checkSignaled())
}

func TestManagedProcHangupNotRunning(t *testing.T) {
	p := managedProc{LoggingCmd: DefaultLoggingCommand("false")}
	assert.ErrorContains(t, p.Hangup(), "not available")
}

func TestManagedProcFailure(t *testing.T) {
	called := false
	onExit := func(err error) {
		called = true
	}

	proc := NewDefaultManagedProc("/xbin/no-such-thing", []string{}, onExit)
	err := proc.Start()
	assert.NonNil(t, err)
	assert.False(t, called)
}
