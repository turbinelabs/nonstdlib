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
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/turbinelabs/test/assert"
)

func makeLogWriter(prefix string) (LogWriter, *bytes.Buffer) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	return NewLogWriter(logger, prefix), &buffer
}

func TestLogWriterReaderFrom(t *testing.T) {
	logWriter, output := makeLogWriter("FROM ")

	s := "abc\nxyzpdq\r\n\r\nstuff\n\ndone\n"
	expected := "FROM abc\nFROM xyzpdq\nFROM \nFROM stuff\nFROM \nFROM done\n"

	reader := strings.NewReader(s)

	_, err := logWriter.ReadFrom(reader)
	assert.Nil(t, err)
	assert.Equal(t, output.String(), expected)
}

func TestLogWriterWrite(t *testing.T) {
	logWriter, output := makeLogWriter("WRITE ")

	logWriter.Write([]byte("xyz"))
	assert.Equal(t, output.String(), "")

	logWriter.Write([]byte("pdq\n"))
	assert.Equal(t, output.String(), "WRITE xyzpdq\n")
	output.Reset()

	logWriter.Write([]byte("test CRLF\r\n"))
	assert.Equal(t, output.String(), "WRITE test CRLF\n")
	output.Reset()

	logWriter.Write([]byte("\n"))
	assert.Equal(t, output.String(), "WRITE \n")
	output.Reset()

	logWriter.Write([]byte("\n\n"))
	assert.Equal(t, output.String(), "WRITE \nWRITE \n")
	output.Reset()

	logWriter.Write([]byte("\r\n"))
	assert.Equal(t, output.String(), "WRITE \n")
	output.Reset()

	logWriter.Write([]byte("\r\n\r\n"))
	assert.Equal(t, output.String(), "WRITE \nWRITE \n")
	output.Reset()

	logWriter.Write([]byte("abc\ndef\n"))
	assert.Equal(t, output.String(), "WRITE abc\nWRITE def\n")
	output.Reset()
}

func TestLogWriterWriteBufferExpansion(t *testing.T) {
	logWriter, output := makeLogWriter("WRITE ")

	logWriter.Write([]byte("abcdefghijklmnopqrstuvwyz\n"))
	output.Reset()

	expected := "WRITE "
	for i := 0; i < 500; i++ {
		x := fmt.Sprintf("(%04d)", i)
		logWriter.Write([]byte(x))
		expected += x
	}
	logWriter.Write([]byte("END\n"))
	expected += "END\n"
	assert.Equal(t, output.String(), expected)
}
