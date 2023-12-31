// nanoid is an opinionated use of the `github.com/matoous/go-nanoid/v2` library.
// It provides default length and helpers to prefix the generated id.
package nanoid

import (
	"fmt"

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
