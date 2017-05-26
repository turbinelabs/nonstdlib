package executor

import "testing"

func TestGoroutineExecTaskDiagnostics(t *testing.T) {
	testExecTaskDiagnostics(t, NewGoroutineExecutor)
}

func TestGoroutineExecAttemptDiagnostics(t *testing.T) {
	testExecAttemptDiagnostics(t, NewGoroutineExecutor)
}
