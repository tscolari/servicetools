package nanoid_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/nanoid"
)

func Test_New(t *testing.T) {
	id := nanoid.New("cus")
	require.Regexp(t, "cus_[a-zA-Z0-9]{17}", id)
}

func Test_Generator(t *testing.T) {
	idFunc := nanoid.Generator("acc")
	id := idFunc()
	require.Regexp(t, "acc_[a-zA-Z0-9]{17}", id)
}
