// Package regexp is a place to put useful regex patterns
package regexp

import "regexp"

var identRegexp = regexp.MustCompile("^([[:alpha:]][[:alnum:]_]*|_[[:alnum:]_]+)$")

// GolangIdentifierRegexp returns a Regexp for validating
// strings as suitable Golang identifiers. It only works for a practical subset
// of possible identifiers, in particular, it does not support non-ascii
// unicode letters and numbers.
func GolangIdentifierRegexp() *regexp.Regexp {
	return identRegexp
}
