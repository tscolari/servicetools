package server

import (
	context "context"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/testhelpers"
	"google.golang.org/grpc"
)

type testWithGRPC struct {
	*WithGRPC
}

func (t *testWithGRPC) ConfigureGRPC(w *WithGRPC) {
	t.WithGRPC = w
}

func Test_WithGRPC(t *testing.T) {
	t.Run("started channel: multiple listeners are allowed", func(t *testing.T) {
		listener1Ok := false
		listener2Ok := false

		withGRPC := NewWithGRPC("localhost:0", grpc.ConnectionTimeout(100*time.Millisecond))
		defer withGRPC.Stop(context.Background())

		go func() {
			<-withGRPC.StartedChan()
			listener1Ok = true
		}()

		go func() {
			<-withGRPC.StartedChan()
			listener2Ok = true
		}()

		testhelpers.Constantly(t, func() bool {
			return (!listener1Ok) && (!listener2Ok)
		}, 250*time.Millisecond, 100*time.Millisecond, "channel listener 1 changed")

		go func() {
			require.NoError(t, withGRPC.Start(context.Background(), slog.Default()))
		}()

		require.Eventually(t, func() bool {
			return listener1Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 1 didn't return")

		require.Eventually(t, func() bool {
			return listener2Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 2 didn't return")

	})

	t.Run("server is reachable", func(t *testing.T) {
		withGRPC := NewWithGRPC("localhost:0", grpc.ConnectionTimeout(100*time.Millisecond))

		go func() {
			require.NoError(t, withGRPC.Start(context.Background(), slog.Default()))
		}()
		defer withGRPC.Stop(context.Background())

		select {
		case <-withGRPC.StartedChan():
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "timed out waiting for server to start")
		}

		conn, err := net.DialTimeout("tcp", withGRPC.address, 100*time.Millisecond)
		require.NoError(t, err)
		defer conn.Close()

		t.Run("all services get mounted", func(t *testing.T) {
			// TODO
		})

		t.Run("closing the server", func(t *testing.T) {
			withGRPC.Stop(context.Background())

			_, err := net.DialTimeout("tcp", withGRPC.address, 50*time.Millisecond)
			require.Error(t, err)
			require.Contains(t, err.Error(), "connection refused")
		})
	})
}
