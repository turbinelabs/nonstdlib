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

package executor

import "testing"

func TestGoroutineExecWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecWithGlobalTimeoutSucceeds(t, NewGoroutineExecutor)
}

func TestGoroutineExecWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecWithGlobalTimeoutTimesOut(t, NewGoroutineExecutor)
}

func TestGoroutineExecWithGlobalTimeoutTimesOutBeforeRetry(t *testing.T) {
	testExecWithGlobalTimeoutTimesOutBeforeRetry(t, NewGoroutineExecutor)
}

func TestGoroutineExecManyWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecManyWithGlobalTimeoutSucceeds(t, NewGoroutineExecutor)
}

func TestGoroutineExecManyWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecManyWithGlobalTimeoutTimesOut(t, NewGoroutineExecutor)
}

func TestGoroutineExecGatheredWithGlobalTimeoutSucceeds(t *testing.T) {
	testExecGatheredWithGlobalTimeoutSucceeds(t, NewGoroutineExecutor)
}

func TestGoroutineExecGatheredWithGlobalTimeoutTimesOut(t *testing.T) {
	testExecGatheredWithGlobalTimeoutTimesOut(t, NewGoroutineExecutor)
}

func TestGoroutineExecWithAttemptTimeoutSucceeds(t *testing.T) {
	testExecWithAttemptTimeoutSucceeds(t, NewGoroutineExecutor)
}

func TestGoroutineExecWithAttemptTimeoutTimesOut(t *testing.T) {
	testExecWithAttemptTimeoutTimesOut(t, NewGoroutineExecutor)
}

func TestGoroutineExecWithGlobalAndAttemptTimeoutsTimesOut(t *testing.T) {
	testExecWithGlobalAndAttemptTimeoutsTimesOut(t, NewGoroutineExecutor)
}
