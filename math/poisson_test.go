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

package math

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/turbinelabs/test/assert"
)

func TestNewPoissonDistribution(t *testing.T) {
	p, err := NewPoissonDistribution(0.0)
	assert.Nil(t, p)
	assert.ErrorContains(t, err, "rate must be greater than 0")

	p, err = NewPoissonDistribution(0.1)
	assert.Nil(t, err)
	assert.NonNil(t, p)
	assert.NonNil(t, p.rng)

	rng := rand.New(rand.NewSource(1))
	p, err = NewPoissonDistributionWithRand(0.1, rng)
	assert.Nil(t, err)
	assert.NonNil(t, p)
	assert.SameInstance(t, p.rng, rng)
}

func TestPoissonDistribution(t *testing.T) {
	p, err := NewPoissonDistribution(1.0)
	assert.Nil(t, err)

	totalDelay := time.Duration(0)
	const n = 1000000
	for i := 0; i < n; i++ {
		totalDelay += p.Next()
	}

	averageDelay := totalDelay / time.Duration(n)
	percentageOfExpected := float64(time.Second-averageDelay) / float64(time.Second) * 100.0

	assert.True(t, math.Abs(percentageOfExpected) < 1.0)
}
