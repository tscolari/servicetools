package validations

import (
	"errors"
	"fmt"

	"github.com/tscolari/servicetools/nanoid"
)

// IsNanoid provides a validation rule that checks if given ID complies with the rules
// defined in the servicetools/nanoid package.
func IsNanoid(prefix, alphabet string, length int) Rule {
	return &idRule{
		prefix:   prefix,
		alphabet: alphabet,
		length:   length,
	}
}

type idRule struct {
	prefix   string
	alphabet string
	length   int
}

func (r *idRule) Validate(value interface{}) error {
	id, ok := value.(string)
	if !ok {
		return fmt.Errorf("the value is not a string, but %t", value)
	}

	if r.prefix == "" {
		if !nanoid.Valid(r.alphabet, r.length, id) {
			return errors.New("the value is not a valid id")
		}
	} else {
		if !nanoid.ValidWithPrefix(r.prefix, r.alphabet, r.length, id) {
			return errors.New("the value is not a valid id")
		}
	}

	return nil
}
