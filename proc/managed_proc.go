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

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"errors"
	"log"
	"syscall"
)

// ManagedProc models a process under management.
type ManagedProc interface {
	// Path returns the path of the command to be run.
	Path() string

	// Args returns a copy of the arguments of the command to be run.
	Args() []string

	// Start the process.
	Start() error

	// Process is running.
	Running() bool

	// Process ran to completion and exited without being signaled.
	Completed() bool

	// Sends SIGHUP to the process. Returns an error if the
	// process is not running.
	Hangup() error

	// Sends SIGQUIT to the process. If the process has already
	// terminated, no error is returned.
	Quit() error

	// Sends SIGKILL to the process. If the process has already
	// terminated, no error is returned.
	Kill() error

	// Sends SIGTERM to the process. If the process has already
	// terminated, no error is returned.
	Term() error

	// Sends SIGUSR1 to the process. If the process has already
	// terminated, no error is returned.
	Usr1() error

	// Waits until the process exits.
	Wait() error
}

type managedProc struct {
	*LoggingCmd

	running bool
	onExit  func(error)
}

// NewDefaultManagedProc constructs a new ManagedProc using
// DefaultLoggingCommand. Invokes onExit (if non-nil) when the process
// stops running. If the process cannot be started, an error is
// returned. In this case onExit is not invoked.
func NewDefaultManagedProc(exe string, args []string, onExit func(error)) ManagedProc {
	return &managedProc{LoggingCmd: DefaultLoggingCommand(exe, args...), onExit: onExit}
}

// NewManagedProc constructs a new ManagedProc with a LoggingCommand
// using the given logger. See NewDefaultManagedProc for details on
// onExit.
func NewManagedProc(
	exe string,
	args []string,
	logger *log.Logger,
	onExit func(error),
) ManagedProc {
	return &managedProc{LoggingCmd: LoggingCommand(logger, exe, args...), onExit: onExit}
}

func (p *managedProc) Path() string {
	return p.LoggingCmd.Path
}

func (p *managedProc) Args() []string {
	args := make([]string, len(p.LoggingCmd.Args))
	copy(args, p.LoggingCmd.Args)
	return args
}

func (p *managedProc) Start() error {
	if err := p.LoggingCmd.Start(); err != nil {
		return err
	}

	go func() {
		err := p.Wait()
		p.running = false
		if p.onExit != nil {
			p.onExit(err)
		}
	}()

	p.running = true

	return nil
}

func (p *managedProc) Running() bool {
	return p.running
}

func (p *managedProc) Completed() bool {
	return p.ProcessState != nil && p.ProcessState.Exited()
}

func (p *managedProc) Hangup() error {
	if p.Process != nil {
		return p.Process.Signal(syscall.SIGHUP)
	}

	return errors.New("process not available")
}

func (p *managedProc) Term() error {
	if p.Process != nil {
		return p.Process.Signal(syscall.SIGTERM)
	}

	return errors.New("process not available")
}

func (p *managedProc) Usr1() error {
	if p.Process != nil {
		return p.Process.Signal(syscall.SIGUSR1)
	}

	return errors.New("process not available")
}
