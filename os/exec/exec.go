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

// Run executes a command then returns (stdout, stderr, error); if the process
// returns success (exit code 0) error is nil.
//
// In all cases the stdout and stderr from the executed command is returned.
func Run(cmd string, args ...string) (string, string, error) {
	execcmd := exec.Command(cmd, args...)

	stdoutBuffer := bytes.Buffer{}
	execcmd.Stdout = &stdoutBuffer

	stderrBuffer := bytes.Buffer{}
	execcmd.Stderr = &stderrBuffer

	err := execcmd.Run()

	return string(stdoutBuffer.Bytes()), string(stderrBuffer.Bytes()), ProcessErr(err)
}

// RunInTerm executes a command redirecting stderr, stdout, and stdin from the
// active TERM.
func RunInTerm(cmd string, args ...string) error {
	exccmd := exec.Command(cmd, args...)
	exccmd.Stdout = os.Stdout
	exccmd.Stdin = os.Stdin
	exccmd.Stderr = os.Stderr
	return ProcessErr(exccmd.Run())
}
