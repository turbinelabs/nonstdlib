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
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestNewExponentialDistribution(t *testing.T) {
	dist := NewExponentialDistribution(0.5)
	assert.Equal(t, dist, ExponentialDistribution(0.5))
}

func TestExponentialDistributionCumulativeProbability(t *testing.T) {
	dist := NewExponentialDistribution(0.5)
	assert.Equal(t, dist.CumulativeProbability(0.0), 0.0)
	assert.Equal(t, dist.CumulativeProbability(-1.0), 0.0)
	assert.Equal(t, dist.CumulativeProbability(0.5), 0.6321205588285577)
}

func TestExponentialDistributionInverseCumulativeProbability(t *testing.T) {
	dist := NewExponentialDistribution(0.5)

	inv, err := dist.InverseCumulativeProbability(0.0)
	assert.Nil(t, err)
	assert.Equal(t, inv, 0.0)

	inv, err = dist.InverseCumulativeProbability(-1.0)
	assert.ErrorContains(t, err, "p must be in the range [0.0, 1.0]")
	assert.True(t, math.IsNaN(inv))

	inv, err = dist.InverseCumulativeProbability(1.1)
	assert.ErrorContains(t, err, "p must be in the range [0.0, 1.0]")
	assert.True(t, math.IsNaN(inv))

	inv, err = dist.InverseCumulativeProbability(1.0)
	assert.Nil(t, err)
	assert.Equal(t, inv, math.Inf(+1))

	inv, err = dist.InverseCumulativeProbability(0.5)
	assert.Nil(t, err)
	assert.Equal(t, inv, 0.34657359027997264)
}
