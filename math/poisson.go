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

package math

import (
	"errors"
	"math/rand"
	"time"
)

// PoissonDistribution models the number of times an event occurs in
// an interval of time. It produces a sequence of time.Duration values
// that represent the delay between the events for an average rate.
type PoissonDistribution struct {
	rng  *rand.Rand
	dist ExponentialDistribution
}

// NewPoissonDistribution creates a new new PoissonDistribution with
// the given average rate. The rate must be greater than 0.
//
// The PoissonDistribution created by this function contains a random
// number source that is not safe for concurrent use (see
// math/rand.NewSource).
func NewPoissonDistribution(rate float64) (*PoissonDistribution, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return NewPoissonDistributionWithRand(rate, rng)
}

// NewPoissonDistributionWithRand creates a new PoissonDistribution
// with the given average rate and random number generator. The rate
// must be greater than 0.
func NewPoissonDistributionWithRand(rate float64, rng *rand.Rand) (*PoissonDistribution, error) {
	if rate <= 0.0 {
		return nil, errors.New("rate must be greater than 0")
	}

	return &PoissonDistribution{
		rng:  rng,
		dist: NewExponentialDistribution(float64(time.Second) / rate),
	}, nil
}

// Next generates the next delay in the sequence.
func (p *PoissonDistribution) Next() time.Duration {
	x, err := p.dist.InverseCumulativeProbability(p.rng.Float64())
	if err != nil {
		// Since math/rand.Float64 produces values inside the
		// valid range for inverseCumulativeProbability, this
		// cannot happen.
		panic(err)
	}
	return time.Duration(x)
}
