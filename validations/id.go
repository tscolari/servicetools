package validations

import (
	"errors"
	"fmt"

	"github.com/tscolari/servicetools/nanoid"
)

func IsID(prefix string) *idRule {
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
