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
	"fmt"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
	testlog "github.com/turbinelabs/test/log"
)

func makeLogWriter(prefix string) (LogWriter, <-chan string) {
	logger, logChannel := testlog.NewChannelLogger(10)
	return NewLogWriter(logger, prefix), logChannel
}

func TestLogWriterReaderFrom(t *testing.T) {
	logWriter, output := makeLogWriter("FROM ")

	s := "abc\nxyzpdq\r\n\r\nstuff\n\ndone\n"

	reader := strings.NewReader(s)

	_, err := logWriter.ReadFrom(reader)
	assert.Nil(t, err)
	assert.Equal(t, <-output, "FROM abc\n")
	assert.Equal(t, <-output, "FROM xyzpdq\n")
	assert.Equal(t, <-output, "FROM \n")
	assert.Equal(t, <-output, "FROM stuff\n")
	assert.Equal(t, <-output, "FROM \n")
	assert.Equal(t, <-output, "FROM done\n")
}

func TestLogWriterWrite(t *testing.T) {
	logWriter, output := makeLogWriter("WRITE ")

	logWriter.Write([]byte("xyz"))
	select {
	case s := <-output:
		assert.Failed(t, fmt.Sprintf("expected no output, got %q", s))
	default:
		// ok
	}

	logWriter.Write([]byte("pdq\n"))
	assert.Equal(t, <-output, "WRITE xyzpdq\n")

	logWriter.Write([]byte("test CRLF\r\n"))
	assert.Equal(t, <-output, "WRITE test CRLF\n")

	logWriter.Write([]byte("\n"))
	assert.Equal(t, <-output, "WRITE \n")

	logWriter.Write([]byte("\n\n"))
	assert.Equal(t, <-output, "WRITE \n")
	assert.Equal(t, <-output, "WRITE \n")

	logWriter.Write([]byte("\r\n"))
	assert.Equal(t, <-output, "WRITE \n")

	logWriter.Write([]byte("\r\n\r\n"))
	assert.Equal(t, <-output, "WRITE \n")
	assert.Equal(t, <-output, "WRITE \n")

	logWriter.Write([]byte("abc\ndef\n"))
	assert.Equal(t, <-output, "WRITE abc\n")
	assert.Equal(t, <-output, "WRITE def\n")
}

func TestLogWriterWriteBufferExpansion(t *testing.T) {
	logWriter, output := makeLogWriter("WRITE ")

	logWriter.Write([]byte("abcdefghijklmnopqrstuvwyz\n"))
	<-output

	expected := "WRITE "
	for i := 0; i < 500; i++ {
		x := fmt.Sprintf("(%04d)", i)
		logWriter.Write([]byte(x))
		expected += x
	}
	logWriter.Write([]byte("END\n"))
	expected += "END\n"
	assert.Equal(t, <-output, expected)
}
