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
	"errors"
	"math"
)

// ExponentialDistribution represents an exponential distribution with
// a given mean.
type ExponentialDistribution float64

// NewExponentialDistribution creates an ExponentialDistribution with
// the given mean.
func NewExponentialDistribution(mean float64) ExponentialDistribution {
	return ExponentialDistribution(mean)
}

// CumulativeProbability returns P(X < x) for this distribution. Based on
// http://mathworld.wolfram.com/ExponentialDistribution.html
func (mean ExponentialDistribution) CumulativeProbability(x float64) float64 {
	if x <= 0.0 {
		return 0.0
	}

	return 1.0 - math.Exp(-x/float64(mean))
}

// InverseCumulativeProbability returns the critical point x, such
// that P(X < x) = p.
func (mean ExponentialDistribution) InverseCumulativeProbability(p float64) (float64, error) {
	if p < 0.0 || p > 1.0 {
		return math.NaN(), errors.New("p must be in the range [0.0, 1.0]")
	}

	if p == 1.0 {
		return math.Inf(+1), nil
	}

	return -float64(mean) * math.Log(1.0-p), nil
}
