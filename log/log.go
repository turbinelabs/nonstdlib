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

// Package log provides infrastructure for topic based logging to files.
//
// Currently all loggers must be specified through the Initialize call and each
// may be accessed through package's Get. Before program termination Close
// should be called to gracefully flush all loggers.
package log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

// Config specifies how a particular logger should operate.
type Config struct {
	// Filename indicates the name of the log file that entries should be sent to.
	Filename string
	// If TeeTo is set entries will also be sent to the indicated Writer.
	TeeTo io.Writer
	// golang log.Logger flags
	Flags int
	// log line prefix
	Prefix string
}

// DefaultConfig provides some reasonable defaults for topic config.
var DefaultConfig = Config{
	Filename: "",
	TeeTo:    os.Stderr,
	Flags:    log.LstdFlags | log.Lshortfile,
	Prefix:   "",
}

// Initialize confgures logging and individual topic based loggers.
//
//   Parameters:
//     stdoutOnly - all loggers will ignore file destinations and produce output
//                  to stdout. Their topic will be prefixed to the log line
//     logRoot - indicates the path directory all log files should be placed in
//     rotateOnSignal - if set to true this package will listen for SIGHUP
//                      signals and reopen the log file they're writing to when
//                      it is received; has no meaning if stdoutOnly is set
//     topicConfig - maps topic to logger config.
//     fileOpenErrCB - in the event a file fails to open during rotation or
//                     initialization this will be called; it will be passed
//                     the topic, the previous file pointer, and the error from
//                     opening the file. If a new file pointer is returned that
//                     will be used as the location logging continues.
//
// If Initialize is called twice it will return an error.
//
// If fileOpenErrCB is nil the default behavior is to set the logger's output
// to Stderr, log the error, and continue operating.
func Initialize(
	stdoutOnly bool,
	logRoot string,
	rotateOnSignal bool,
	topicConfig map[string]Config,
	fileOpenErrCB func(string, io.WriteCloser, error) io.WriteCloser,
) error {
	return state.Initialize(stdoutOnly, logRoot, rotateOnSignal, topicConfig, fileOpenErrCB)
}

func (s *loggerstate) Initialize(
	stdoutOnly bool,
	logRoot string,
	rotateOnSignal bool,
	topicConfig map[string]Config,
	fileOpenErrCB func(string, io.WriteCloser, error) io.WriteCloser,
) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.initialized {
		return errors.New("Attempting to initialize log infrastructure twice")
	}

	s.observeSignal = rotateOnSignal
	s.stdoutOnly = stdoutOnly
	s.logRoot = logRoot
	s.configs = topicConfig
	if fileOpenErrCB != nil {
		s.fileErrCB = fileOpenErrCB
	} else {
		s.fileErrCB = defaultFileErrCB
	}

	for tpc, cfg := range topicConfig {
		pfix := cfg.Prefix
		if stdoutOnly {
			pfix = string(tpc) + ":" + pfix
		}

		// even if not stdout only we initialize the logger with stdout because
		// reopenfile will handle setting the logger's output io.Writer
		logger := log.New(os.Stdout, pfix, cfg.Flags)

		s.loggers[tpc] = logger
		if !stdoutOnly {
			s.reopenFile(tpc, cfg, logger)
		}
	}

	if !s.stdoutOnly && rotateOnSignal {
		signalChannel := make(chan os.Signal, 1)
		go handleHup(signalChannel)
		signal.Notify(signalChannel, syscall.SIGHUP)
	}

	s.initialized = true
	return nil
}

func (s *loggerstate) Close() {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if state.observeSignal {
		signal.Reset(syscall.SIGHUP)
	}

	for t, f := range state.openFiles {
		// stop tracking this file
		delete(state.openFiles, t)

		// and stop writing to it
		if l := state.loggers[t]; l != nil {
			l.SetOutput(os.Stderr)
			l.SetPrefix("[closed] " + l.Prefix())
		}

		f.Close()
	}

	state.closed = true
}

// Close will Close on all files that have been opened for logging.
//
// Accessing a topic-specific logger via Get after closing them will cause
// output to be logged to stderr instead of the desired file but will not panic.
//
// Log entries that were previously logging to files will also have a [closed]
// prefix added to indicate their underlying Writer has been closed and changed.
//
// Signals for logrotation will no longer be observed.
func Close() {
	state.Close()
}

// A default handler that just logs the error and panics instead of trying
// to recover or continue.
func LogAndPanic(topic string, prev io.WriteCloser, err error) io.WriteCloser {
	log.Panicf("Could not open %s: %s", topic, err)
	return nil
}

type loggerstate struct {
	mutex         sync.RWMutex
	initialized   bool
	closed        bool
	stdoutOnly    bool
	observeSignal bool
	fileErrCB     func(string, io.WriteCloser, error) io.WriteCloser

	logRoot   string
	configs   map[string]Config
	openFiles map[string]io.WriteCloser
	loggers   map[string]*log.Logger

	// provided as function pointer to make testing possible
	newLogFile func(string, string) (io.WriteCloser, error)
}

var state = loggerstate{
	configs:    make(map[string]Config),
	loggers:    make(map[string]*log.Logger),
	openFiles:  make(map[string]io.WriteCloser),
	mutex:      sync.RWMutex{},
	newLogFile: newLogFile,
}

// construct path to a log file and and return a file pointer
func newLogFile(root, file string) (io.WriteCloser, error) {
	path := filepath.Join(root, file)
	logfile, err := os.OpenFile(
		path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|0660)

	if err != nil {
		return nil, err
	}

	return logfile, nil
}

// Opens a file of the appropriate name and sets a writer to it as output
// for a logger. If there is an existing WriteCloser saved for a given topic
// this will also close it.
//
// On failure to reopen will invoke the configured callback allowing the
// user a chance to recover. If the recovered io.WriteCloser is still nil after
// the callback the loger will be directed at os.Stderr with an updated
// prefix of: "[in error: $topic] $originalPrefix"
//
// If the logger is set to operate on stdout only no action will be taken
func (ls *loggerstate) reopenFile(
	topic string,
	cfg Config,
	logger *log.Logger,
) {
	// don't need to take any action -- only logging to stdout so there is no file
	// file to reopen
	if ls.stdoutOnly {
		return
	}

	// save the old file to close
	oldWriteCloser, hasOldFile := ls.openFiles[topic]
	delete(ls.openFiles, topic)

	filename := string(topic) + ".log"
	if cfg.Filename != "" {
		filename = cfg.Filename
	}

	var writer io.Writer
	var destWriter io.WriteCloser

	// get a pointer to a new logfile
	destWriter, err := ls.newLogFile(ls.logRoot, filename)

	if err != nil {
		// give the client a chance to decide how to recover
		destWriter = state.fileErrCB(topic, oldWriteCloser, err)
	}
	writer = destWriter

	// if configured for teeing output;
	//   and the destWriter is non-nil (i.e. either recovered or no error)
	if cfg.TeeTo != nil && destWriter != nil {
		// then we create a multiwriter pointing at the tee target
		writer = io.MultiWriter(destWriter, cfg.TeeTo)
	}

	// even in the case that we got the old writer back we need to reinsert
	// it since we previously removed it
	if destWriter != nil {
		ls.openFiles[topic] = destWriter
	}

	// if we don't have a valid writer by now fail to stderr and update prefix
	// to indicate original source
	if writer == nil {
		writer = os.Stderr
		logger.SetPrefix(fmt.Sprintf("[in error: %s] %s", topic, logger.Prefix()))
	}

	logger.SetOutput(writer)

	if hasOldFile && destWriter != oldWriteCloser {
		oldWriteCloser.Close()
	}
}

func handleHup(channel chan os.Signal) {
	for {
		<-channel
		state.mutex.Lock()

		if state.closed {
			state.mutex.Unlock()
			return
		}

		for tpc, cfg := range state.configs {
			logger := state.loggers[tpc]
			if logger != nil {
				state.reopenFile(tpc, cfg, logger)
			}
		}

		// can't defer because we don't exit this func
		state.mutex.Unlock()
	}
}

type lw struct {
	*log.Logger
}

func (l lw) Write(data []byte) (n int, err error) {
	l.Print(string(data))
	return len(data), nil
}

// ToWriter provides an io.Writer adapter that converts Writes into log actions.
func ToWriter(logger *log.Logger) io.Writer {
	return lw{logger}
}

// Failsafe is used to prevent dataloss in the case where a log client attempts
// to access a log before it has been initialized.
var Failsafe *log.Logger = log.New(os.Stderr, "failsafe: ", log.LstdFlags|log.Llongfile)

func topicalFailsafe(t string) *log.Logger {
	return log.New(
		os.Stderr, fmt.Sprintf("failsafe - %s: ", t), log.LstdFlags|log.Llongfile)
}

func (s *loggerstate) Get(topic string) *log.Logger {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.initialized {
		Failsafe.Printf("Attempting to access %s loggers before initialization", topic)
		return Failsafe
	}

	if s.closed {
		Failsafe.Printf("Attempting to access %s logger after Close() call", topic)
		// we let this fall through because we only clean up the file descriptor
		// not the logger in Close()
	}

	if s.loggers[topic] == nil {
		return topicalFailsafe(topic)
	}

	return s.loggers[topic]
}

// Get provides access to a configured logger. If a logger for the requested
// topic has not been configured return Failsafe and log the issue.
//
// NOTE: use of Get will acquire a read lock on the state mutex. This by itself
// isn't problematic but when we are handling a SIGHUP we also acquire a lock
// while we reopen log files. If Get usage is gonig to be frequent spend some
// time to minimize that critical section.
func Get(topic string) *log.Logger {
	return state.Get(topic)
}

func defaultFileErrCB(topic string, prev io.WriteCloser, err error) io.WriteCloser {
	Failsafe.Printf("Could not open new file for %s: %s", topic, err)
	return nil
}
