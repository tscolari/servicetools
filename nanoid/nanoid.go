package nanoid

import (
	"fmt"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	// maxAttempts defines how many times the helper function will try to generate
	// an ID before panicking - in case the underlying nanoid package fails to generate.
	maxAttempts = 100
)

// New will return a new nanoid with the given alphabet and length.
// The underlying library might fail to generate an id, and this
// wrapper function will try to generate it for `maxAttempts`, but
// in the unlikely event of it not succeeding in any of them it will panic.
func New(alphabet string, length int) string {
	var err error

	for i := 0; i < maxAttempts; i++ {
		var id string
		id, err = gonanoid.Generate(alphabet, length)
		if err == nil {
			return id
		}
	}

	if err != nil {
		panic("nanoid generation failed: " + err.Error())
	}

	return ""
}

// Generator returns a function that generates nanoids with the given configuration
func Generator(alphabet string, length int) func() string {
	return func() string {
		return New(alphabet, length)
	}
}

// Valid verifies if the given nanoid is valid using the given rules.
func Valid(alphabet string, length int, nanoid string) bool {
	m, _ := regexp.Match(fmt.Sprintf("^[%s]{%d}$", alphabet, length), []byte(nanoid))
	return m
}
