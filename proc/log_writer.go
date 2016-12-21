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
	"bufio"
	"io"
	"log"
	"os"
)

type LogWriterState struct {
	buffer []byte
	pos    int
	end    int
	cap    int
}

// Copies data written to pipe (normally from an os.Process) to the
// given log (or the default log).
type LogWriter struct {
	prefix string
	logger *log.Logger
	state  *LogWriterState
}

var _ io.WriteCloser = LogWriter{}
var _ io.ReaderFrom = LogWriter{}

// Creates a new LogWriter that emits lines to the given logger with
// the given prefix.
func NewLogWriter(logger *log.Logger, prefix string) LogWriter {
	return LogWriter{prefix: prefix, logger: logger, state: &LogWriterState{}}
}

// Creates a new LogWriter that emits lines to the default logger
// (on stderr) with the given prefix.
func NewDefaultLogWriter(prefix string) LogWriter {
	logger := log.New(os.Stderr, log.Prefix(), log.Flags())
	return NewLogWriter(logger, prefix)
}

func (w LogWriter) log(line string) {
	w.logger.Printf("%s%s", w.prefix, line)
}

func (w LogWriter) Close() error {
	s := w.state
	if s.buffer != nil && s.end > s.pos {
		w.log(string(s.buffer[s.pos:s.end]))
	}
	return nil
}

func (w LogWriter) Write(b []byte) (int, error) {
	s := w.state

	// First call to Write?
	if s.buffer == nil {
		s.pos = 0
		s.end = 0
		s.cap = 2048
		s.buffer = make([]byte, s.cap)
	}

	// If there is insufficient room for b at the end of the
	// buffer, shift contents of buffer to the start of the buffer
	num := len(b)
	if s.pos > 0 && s.cap-s.end < num {
		copy(s.buffer, s.buffer[s.pos:s.end])
		s.end -= s.pos
		s.pos = 0
	}

	// If there is insufficient room for b at the end of the
	// buffer, expand the buffer until there is.
	for s.cap-s.end < num {
		s.cap *= 2
		newBuffer := make([]byte, s.cap)
		copy(newBuffer, s.buffer[s.pos:s.end])
		s.buffer = newBuffer
	}

	copy(s.buffer[s.end:], b)
	s.end += num

	// Iterate over the live slice of the buffer and emit any
	// complete lines (where lines end with LF or CRLF).
	startPos := s.pos
	slice := s.buffer[s.pos:s.end]
	for i, b := range slice {
		if b == '\n' {
			// convert s.pos to slice index
			first := s.pos - startPos
			last := i
			if i > 0 && slice[i-1] == '\r' {
				last = i - 1
			}
			w.log(string(slice[first:last]))

			// update live slice start to point just after the LF
			s.pos = startPos + i + 1
		}
	}

	return num, nil
}

func (w LogWriter) ReadFrom(r io.Reader) (int64, error) {
	var read int64
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		w.log(text)
		read += int64(len(text)) + 1
	}
	err := scanner.Err()
	return read, err
}
