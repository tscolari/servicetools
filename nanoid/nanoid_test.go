package nanoid_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/nanoid"
)

func Test_New(t *testing.T) {
	id := nanoid.New(nanoid.AlphabetDefault, 15)
	require.Regexp(t, "[a-zA-Z0-9]{15}", id)

	id = nanoid.New(nanoid.AlphabetNumeric, 5)
	require.Regexp(t, "[0-9]{5}", id)
}

func Test_Generator(t *testing.T) {
	idFunc := nanoid.Generator(nanoid.AlphabetDefault, 15)
	require.Regexp(t, "[a-zA-Z0-9]{15}", idFunc())
	require.Regexp(t, "[a-zA-Z0-9]{15}", idFunc())
	require.Regexp(t, "[a-zA-Z0-9]{15}", idFunc())

	idFunc = nanoid.Generator(nanoid.AlphabetLowerCase, 4)
	require.Regexp(t, "[a-z]{4}", idFunc())
	require.Regexp(t, "[a-z]{4}", idFunc())
	require.Regexp(t, "[a-z]{4}", idFunc())
}

func Test_Valid(t *testing.T) {
	id := nanoid.New(nanoid.AlphabetDefault, 15)
	require.True(t, nanoid.Valid(nanoid.AlphabetDefault, 15, id))

	id = nanoid.New(nanoid.AlphabetDefault, 15)
	require.False(t, nanoid.Valid(nanoid.AlphabetDefault, 15, id[:13]))

	id = nanoid.New(nanoid.AlphabetDefault, 15) + "a"
	require.False(t, nanoid.Valid(nanoid.AlphabetDefault, 15, id))
}
