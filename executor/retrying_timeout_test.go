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

package executor

import "testing"

func TestRetryingExecExecWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecWithGlobalTimeoutSucceeds(t, NewRetryingExecutor)
}

func TestRetryingExecExecWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecWithGlobalTimeoutTimesOut(t, NewRetryingExecutor)
}

func TestRetryingExecExecWithGlobalTimeoutTimesOutBeforeRetry(t *testing.T) {
	testExecWithGlobalTimeoutTimesOutBeforeRetry(t, NewRetryingExecutor)
}

func TestRetryingExecExecManyWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecManyWithGlobalTimeoutSucceeds(t, NewRetryingExecutor)
}

func TestRetryingExecExecManyWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecManyWithGlobalTimeoutTimesOut(t, NewRetryingExecutor)
}

func TestRetryingExecExecGatheredWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecGatheredWithGlobalTimeoutSucceeds(t, NewRetryingExecutor)
}

func TestRetryingExecExecGatheredWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecGatheredWithGlobalTimeoutTimesOut(t, NewRetryingExecutor)
}

func TestRetryingExecExecWithAttemptTimeoutSucceeds(t *testing.T) {
	testExecWithAttemptTimeoutSucceeds(t, NewRetryingExecutor)
}

func TestRetryingExecExecWithAttemptTimeoutTimesOut(t *testing.T) {
	testExecWithAttemptTimeoutTimesOut(t, NewRetryingExecutor)
}

func TestRetryingExecExecWithGlobalAndAttemptTimeoutsTimesOut(t *testing.T) {
	testExecWithGlobalAndAttemptTimeoutsTimesOut(t, NewRetryingExecutor)
}
