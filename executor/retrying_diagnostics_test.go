package executor

import "testing"

func TestRetryingExecTaskDiagnostics(t *testing.T) {
	testExecTaskDiagnostics(t, NewRetryingExecutor)
}

func TestRetryingExecAttemptDiagnostics(t *testing.T) {
	testExecAttemptDiagnostics(t, NewRetryingExecutor)
}
