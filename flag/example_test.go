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
)

func ExampleNewStrings() {
	var strings Strings // typically a field in a struct
	strings = NewStrings()

	flagset := flag.NewFlagSet("example", flag.PanicOnError)
	flagset.Var(
		&strings,
		"x",
		"Flag help.",
	)

	flagset.Parse([]string{"-x=a,b,c"})

	for _, selected := range strings.Strings {
		fmt.Println(selected)
	}
	// Output:
	// a
	// b
	// c
}

func ExampleNewStringsWithConstraint() {
	var strings Strings // typically a field in a struct
	strings = NewStringsWithConstraint("choice1", "option2", "possibility3")

	flagset := flag.NewFlagSet("example", flag.PanicOnError)
	flagset.Var(
		&strings,
		"x",
		"Flag help. Allowed values: "+strings.ValidValuesDescription(),
	)

	flagset.Parse([]string{"-x=choice1,possibility3"})

	for _, selected := range strings.Strings {
		fmt.Println(selected)
	}
	// Output:
	// choice1
	// possibility3
}

func ExampleStrings_withDelimiter() {
	var strings Strings // typically a field in a struct
	strings = Strings{Delimiter: ";"}

	flagset := flag.NewFlagSet("example", flag.PanicOnError)
	flagset.Var(
		&strings,
		"x",
		"Flag help. Allowed values: "+strings.ValidValuesDescription(),
	)

	flagset.Parse([]string{"-x=one;two"})

	for _, selected := range strings.Strings {
		fmt.Println(selected)
	}
	// Output:
	// one
	// two
}

func ExampleNewChoice() {
	var choice Choice // typically a field in a struct
	choice = NewChoice("a", "b", "c")

	flagset := flag.NewFlagSet("example", flag.PanicOnError)
	flagset.Var(
		&choice,
		"x",
		"Flag help. Allowed values: "+choice.ValidValuesDescription(),
	)

	flagset.Parse([]string{"-x=c"})

	fmt.Println(choice.String())
	// Output:
	// c
}
