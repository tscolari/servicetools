package validations

import (
	"errors"
	"fmt"

	"github.com/tscolari/servicetools/nanoid"
)

// IsID provides a validation rule that checks if given ID complies with the rules
// defined in the servicetools/nanoid package.
func IsID(prefix string) Rule {
	return &idRule{
		prefix: prefix,
	}
}

type idRule struct {
	prefix string
}

func (r *idRule) Validate(value interface{}) error {
	id, ok := value.(string)
	if !ok {
		return fmt.Errorf("the value is not a string, but %t", value)
	}

	if !nanoid.Valid(r.prefix, id) {
		return errors.New("the value is not a valid id")
	}

	return nil
}
