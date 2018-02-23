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

package time

import (
	"context"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestDefaultSource(t *testing.T) {
	source := NewSource()

	for i := 0; i < 10; i++ {
		before := time.Now()
		now := source.Now()
		after := time.Now()

		assert.True(t, before.Before(now) || before.Equal(now))
		assert.True(t, after.After(now) || after.Equal(now))
	}

	delay := 1 * time.Millisecond

	timer := source.NewTimer(delay)
	defer timer.Stop()
	<-timer.C()

	ctxt1, cancel1 := source.NewContextWithDeadline(context.TODO(), source.Now().Add(delay))
	defer cancel1()
	<-ctxt1.Done()

	ctxt2, cancel2 := source.NewContextWithTimeout(context.TODO(), delay)
	defer cancel2()
	<-ctxt2.Done()
}
