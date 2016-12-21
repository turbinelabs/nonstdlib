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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	Ready    = "READY"
	Done     = "DONE"
	Signaled = "SIGNALED"

	ReceivedFmt = "RECV'D: %v\n"
)

// proc-test is a trivial application that responds to signals such as those
// sent by proc.LoggingCommand and proc.ManagedProc

func usage() {
	fmt.Printf("usage: %s <delay-in-seconds>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	delay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR1)

	fmt.Println(Ready)
	timeout := time.After(time.Duration(delay) * time.Second)
	select {
	case s := <-signals:
		fmt.Printf(ReceivedFmt, s)
		fmt.Println(Signaled)
	case <-timeout:
		fmt.Println(Done)
	}
}
