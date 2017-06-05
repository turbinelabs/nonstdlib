package strings

import (
	"strings"
)

// PadLeft prepends string s with i spaces.
func PadLeft(s string, i int) string {
	return PadLeftWith(s, i, " ")
}

// PadLeftWith prepends string s with i instances of padding p.
func PadLeftWith(s string, i int, p string) string {
	if i < 0 {
		return s
	}
	padding := strings.Repeat(p, i)

	ss := strings.Split(s, "\n")
	for i, e := range ss {
		ss[i] = padding + e
	}

	return strings.Join(ss, "\n")
}
