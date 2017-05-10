// +build baketest

// Package main runs the github.com/turbinelabs/nonstdlib/executor/test.BakeTestCLI()
// function, but is only built if the build tag "baketest" is set.
package main

import "github.com/turbinelabs/nonstdlib/executor/test"

func main() {
	test.BakeTestCLI()
}
