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

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"time"
)

// Timer is a trivial wrapper around time.Timer which allows a timer
// to be mocked and/or replaced with an implementation that can be
// triggered deterministically.
type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

// NewTimer creates a new Timer, wrapping a time.Timer that will
// expire after the given duration.
func NewTimer(d time.Duration) Timer {
	return &timer{time.NewTimer(d)}
}

// AfterFunc creates a new Timer, wrapping a time.Timer that will
// invoke the given function after the given duration.
func AfterFunc(d time.Duration, f func()) Timer {
	return &timer{time.AfterFunc(d, f)}
}

type timer struct {
	*time.Timer
}

var _ Timer = &timer{}

func (t *timer) C() <-chan time.Time {
	return t.Timer.C
}
