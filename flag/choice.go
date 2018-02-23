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

package flag

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/turbinelabs/nonstdlib/arrays/indexof"
)

// Choice conforms to the flag.Value and flag.Getter interfaces, and
// can be used populate a slice of strings from a flag.Flag.
type Choice struct {
	// Populated from the command line.
	Choice *string

	// All possible values allowed to appear in Choice.
	AllowedValues []string
}

var _ flag.Getter = &Choice{}
var _ ConstrainedValue = &Choice{}

// NewChoice produces a Choice with a set of allowed values.
func NewChoice(allowedValues ...string) Choice {
	return Choice{AllowedValues: allowedValues}
}

// WithDefault assigns a default value. If value is not a valid choice
// it is ignored.
func (cv Choice) WithDefault(value string) Choice {
	cv.Set(value)
	return cv
}

func allowedValuesToDescription(allowedValues []string) string {
	quoted := make([]string, len(allowedValues))
	for i, v := range allowedValues {
		quoted[i] = strconv.Quote(v)
	}

	switch n := len(quoted); n {
	case 0:
		return ""
	case 1:
		return quoted[0]
	case 2:
		return strings.Join(quoted, " or ")
	default:
		commaSep := strings.Join(quoted[0:n-1], ", ")
		return commaSep + ", or " + quoted[n-1]
	}
}

// ValidValuesDescription returns a string describing the allowed
// values for this Choice. For example: "a", "b", or "c".
func (cv *Choice) ValidValuesDescription() string {
	return allowedValuesToDescription(cv.AllowedValues)
}

// String returns the current value of the Choice.
func (cv *Choice) String() string {
	if cv.Choice != nil {
		return *cv.Choice
	}
	return ""
}

// Set sets the current value of the Choice, returning an error if the
// value is not one of the available choices.
func (cv *Choice) Set(value string) error {
	if indexof.String(cv.AllowedValues, value) == indexof.NotFound {
		return fmt.Errorf(
			"invalid flag value: %s, must be one of %s",
			value,
			strings.Join(cv.AllowedValues, ", "),
		)
	}

	cv.Choice = &value
	return nil
}

// Get retrieves the current value of the Choice as an interface{}.
func (cv *Choice) Get() interface{} {
	return cv.Choice
}
