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

func TestRetryingExecExecMany(t *testing.T) {
	testExecMany(t, NewRetryingExecutor)
}

func TestRetryingExecExecManyWithFailure(t *testing.T) {
	testExecManyWithFailure(t, NewRetryingExecutor)
}

func TestRetryingExecExecManyNoop(t *testing.T) {
	testExecManyNoop(t, NewRetryingExecutor)
}

func TestRetryingExecExecGathered(t *testing.T) {
	testExecGathered(t, NewRetryingExecutor)
}

func TestRetryingExecExecGatheredWithError(t *testing.T) {
	testExecGatheredWithError(t, NewRetryingExecutor)
}

func TestRetryingExecExecGatheredNoCallback(t *testing.T) {
	testExecGatheredNoCallback(t, NewRetryingExecutor)
}

func TestRetryingExecExecGatheredNoop(t *testing.T) {
	testExecGatheredNoop(t, NewRetryingExecutor)
}
