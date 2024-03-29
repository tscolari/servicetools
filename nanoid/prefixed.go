// Package nanoid is an opinionated use of the `github.com/matoous/go-nanoid/v2` library.
// It provides default length and helpers to prefix the generated id.
package nanoid

import (
	"fmt"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	length   = 17
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVXZWKYabcdefghijklmnopqrstuvxzwky0123456789"
)

// New returns a new nanoid ID using the given prefix and preset length and alphabet.
func New(prefix string) string {
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic("nanoid generation failed: " + err.Error())
	}

	return fmt.Sprintf("%s_%s", prefix, id)
}

// Generator returns a function that generates nanoids with the given prefix.
func Generator(prefix string) func() string {
	return func() string {
		return New(prefix)
	}
}

// Valid verifies if the given nanoid is valid using the
// rules of this package and the given prefix.
func Valid(prefix, nanoid string) bool {
	m, _ := regexp.Match(fmt.Sprintf("^%s_[%s]{%d}$", prefix, alphabet, length), []byte(nanoid))
	return m
}
