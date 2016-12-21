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
)

const (
	procTest = "proc-test"
)

var (
	procTestReady    = fmt.Sprintf("[%s stdout] READY\n", procTest)
	procTestComplete = fmt.Sprintf("[%s stdout] DONE\n", procTest)
	procTestSignaled = fmt.Sprintf("[%s stdout] SIGNALED\n", procTest)
)

type procTestOutput struct {
	*bytes.Buffer
}

func procTestAndArgs(seconds int) (string, []string) {
	return procTest, []string{fmt.Sprintf("%d", seconds)}
}

func makeProcTestOutput() (*log.Logger, *procTestOutput) {
	output := &procTestOutput{&bytes.Buffer{}}
	return log.New(output, "", 0), output
}

func (o *procTestOutput) checkReady() bool {
	return strings.HasPrefix(o.String(), procTestReady)
}

func (o *procTestOutput) checkCompleted() bool {
	return strings.HasSuffix(o.String(), procTestComplete)
}

func (o *procTestOutput) checkSignaled() bool {
	return strings.HasSuffix(o.String(), procTestSignaled)
}
