package nanoid

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	length   = 17
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVXZWKYabcdefghijklmnopqrstuvxzwky0123456789"
)

func New(prefix string) string {
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic("nanoid generation failed: " + err.Error())
	}

	return fmt.Sprintf("%s_%s", prefix, id)
}

func Generator(prefix string) func() string {
	return func() string {
		return New(prefix)
	}
}
