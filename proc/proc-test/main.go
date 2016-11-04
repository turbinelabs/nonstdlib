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
