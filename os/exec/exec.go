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

// Package exec provides extensions of the os/exec package that provides
// streamlined execution of a command.
package exec

import (
	"bytes"
	"os"
	"os/exec"
)

// ProcessErr converts the error value from an exec.Command execution into nil
// if the command exited with success as determined by checking err.Success().
func ProcessErr(e error) error {
	switch t := e.(type) {
	case *exec.ExitError:
		if !t.Success() {
			return e
		}
	default:
		if e != nil {
			return e
		}
	}

	return nil
}

// RunCmd executes a command then returns (stdout, stderr, error); if
// the process returns success (exit code 0) error is nil.
//
// In all cases the stdout and stderr from the executed command is returned.
func RunCmd(cmd *exec.Cmd) (string, string, error) {
	stdoutBuffer := bytes.Buffer{}
	cmd.Stdout = &stdoutBuffer

	stderrBuffer := bytes.Buffer{}
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()

	return string(stdoutBuffer.Bytes()), string(stderrBuffer.Bytes()), ProcessErr(err)
}

// Run executes a command constructed from the string arguments, then returns
// (stdout, stderr, error); if the process returns success (exit code 0) error
// is nil.
//
// In all cases the stdout and stderr from the executed command is returned.
func Run(cmd string, args ...string) (string, string, error) {
	return RunCmd(exec.Command(cmd, args...))
}

// RunCmdInTerm executes a command redirecting stderr, stdout, and
// stdin from the active TERM.
func RunCmdInTerm(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return ProcessErr(cmd.Run())
}

// RunInTerm executes a command constructed from the string arguments,
// redirecting stderr, stdout, and stdin from the active TERM.
func RunInTerm(cmd string, args ...string) error {
	return RunCmdInTerm(exec.Command(cmd, args...))
}
