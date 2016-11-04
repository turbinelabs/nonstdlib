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
