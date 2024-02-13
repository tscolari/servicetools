package testhelpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
