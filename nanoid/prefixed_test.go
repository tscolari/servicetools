package nanoid_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/nanoid"
)

func Test_NewPrefixed(t *testing.T) {
	id := nanoid.NewPrefixed("cus", nanoid.AlphabetDefault, 15)
	require.Regexp(t, "cus_[a-zA-Z0-9]{15}", id)

	id = nanoid.NewPrefixed("cus", nanoid.AlphabetNumeric, 5)
	require.Regexp(t, "cus_[0-9]{5}", id)
}

func Test_PrefixedGenerator(t *testing.T) {
	idFunc := nanoid.PrefixedGenerator("acc", nanoid.AlphabetDefault, 15)
	require.Regexp(t, "acc_[a-zA-Z0-9]{15}", idFunc())
	require.Regexp(t, "acc_[a-zA-Z0-9]{15}", idFunc())
	require.Regexp(t, "acc_[a-zA-Z0-9]{15}", idFunc())

	idFunc = nanoid.PrefixedGenerator("acc", nanoid.AlphabetLowerCase, 4)
	require.Regexp(t, "acc_[a-z]{4}", idFunc())
	require.Regexp(t, "acc_[a-z]{4}", idFunc())
	require.Regexp(t, "acc_[a-z]{4}", idFunc())
}

func Test_ValidWithPrefix(t *testing.T) {
	id := nanoid.NewPrefixed("cus", nanoid.AlphabetDefault, 15)
	require.True(t, nanoid.ValidWithPrefix("cus", nanoid.AlphabetDefault, 15, id))

	id = nanoid.NewPrefixed("nop", nanoid.AlphabetDefault, 15)
	require.False(t, nanoid.ValidWithPrefix("cus", nanoid.AlphabetDefault, 15, id))

	id = nanoid.NewPrefixed("cus", nanoid.AlphabetDefault, 15)
	require.False(t, nanoid.ValidWithPrefix("cus", nanoid.AlphabetDefault, 15, id[:14]))

	id = nanoid.NewPrefixed("cus", nanoid.AlphabetDefault, 15) + "a"
	require.False(t, nanoid.ValidWithPrefix("cus", nanoid.AlphabetDefault, 15, id))

	id = nanoid.NewPrefixed("acus", nanoid.AlphabetDefault, 15)
	require.False(t, nanoid.ValidWithPrefix("cus", nanoid.AlphabetDefault, 15, id))
}
