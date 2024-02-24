package testhelpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Constantly will execute the assertion for `waitFor` period of time at every `tick` duration.
// For the test to pass, all calls to assertion must return true within that period.
func Constantly(t *testing.T, assertion func() bool, waitFor, tick time.Duration, msgs ...interface{}) {
	timeout := time.NewTimer(waitFor)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			return
		case <-time.After(tick):
			require.True(t, assertion(), msgs...)
		}
	}
}
