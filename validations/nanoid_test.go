package validations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/validations"
)

func Test_IsNanoid(t *testing.T) {

	t.Run("With prefix", func(t *testing.T) {
		rule := validations.IsNanoid("abc", "1234567890", 5)

		testCases := []struct {
			id    string
			valid bool
		}{
			{id: "abc_02312", valid: true},
			{id: "abC_02312", valid: false},
			{id: "abc_023123", valid: false},
			{id: "abc_0231", valid: false},
			{id: "abd_02312", valid: false},
			{id: "abc_0231a", valid: false},
			{id: "02345", valid: false},
		}

		for _, tc := range testCases {
			t.Run(tc.id, func(t *testing.T) {
				require.Equal(t, tc.valid, rule.Validate(tc.id) == nil)
			})
		}
	})

	t.Run("Without prefix", func(t *testing.T) {
		rule := validations.IsNanoid("", "abcdefg", 5)

		testCases := []struct {
			id    string
			valid bool
		}{
			{id: "abcde", valid: true},
			{id: "Abcde", valid: false},
			{id: "abcdef", valid: false},
			{id: "abcd", valid: false},
			{id: "abc_abcdef", valid: false},
			{id: "abcdz", valid: false},
			{id: "abcd0", valid: false},
		}

		for _, tc := range testCases {
			t.Run(tc.id, func(t *testing.T) {
				require.Equal(t, tc.valid, rule.Validate(tc.id) == nil)
			})
		}
	})
}
