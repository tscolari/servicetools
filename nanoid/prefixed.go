// Package nanoid is an opinionated use of the `github.com/matoous/go-nanoid/v2` library.
// It provides default length and helpers to prefix the generated id.
package nanoid

import (
	"fmt"
	"regexp"
)

const (
	AlphabetDefault   = "ABCDEFGHIJKLMNOPQRSTUVXZWKYZabcdefghijklmnopqrstuvxzwkyz0123456789"
	AlphabetNumeric   = "0123456789"
	AlphabetLowerCase = "abcdefghijklmnopqrstuvxzwkyz"
	AlphabetUpperCase = "ABCDEFGHIJKLMNOPQRSTUVXZWKYZ"
)

// NewPrefixed returns a new nanoid ID using the given prefix, length and alphabet.
func NewPrefixed(prefix, alphabet string, length int) string {
	id := New(alphabet, length)
	return fmt.Sprintf("%s_%s", prefix, id)
}

// PrefixedGenerator returns a function that generates nanoids with the given prefix.
func PrefixedGenerator(prefix, alphabet string, length int) func() string {
	return func() string {
		return NewPrefixed(prefix, alphabet, length)
	}
}

// ValidWithPrefixWithPrefix verifies if the given nanoid is valid using the
// rules of this package and the given prefix.
func ValidWithPrefix(prefix, alphabet string, length int, nanoid string) bool {
	m, _ := regexp.Match(fmt.Sprintf("^%s_[%s]{%d}$", prefix, alphabet, length), []byte(nanoid))
	return m
}
