package server

import (
	context "context"
	"fmt"
	slog "log/slog"
	"net"
	http "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/testhelpers"
)

func Test_WithHTTP(t *testing.T) {
	t.Run("started channel: multiple listeners are allowed", func(t *testing.T) {
		listener1Ok := false
		listener2Ok := false

		withHTTP := NewWithHTTP("localhost:0")

		go func() {
			<-withHTTP.StartedChan()
			listener1Ok = true
		}()

		go func() {
			<-withHTTP.StartedChan()
			listener2Ok = true
		}()

		testhelpers.Constantly(t, func() bool {
			return (!listener1Ok) && (!listener2Ok)
		}, 300*time.Millisecond, 100*time.Millisecond, "channel listener changed")

		go func() {
			require.NoError(t, withHTTP.Start(context.Background(), slog.Default()))
		}()

		require.Eventually(t, func() bool {
			return listener1Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 1 didn't return")

		require.Eventually(t, func() bool {
			return listener2Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 2 didn't return")

	})

	t.Run("server is reachable", func(t *testing.T) {
		withHTTP := NewWithHTTP("localhost:0")

		go func() {
			require.NoError(t, withHTTP.Start(context.Background(), slog.Default()))
		}()

		select {
		case <-withHTTP.StartedChan():
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "timed out waiting for server to start")
		}

		conn, err := net.DialTimeout("tcp", withHTTP.address, 100*time.Millisecond)
		require.NoError(t, err)
		require.NoError(t, conn.Close())

		t.Run("closing the server", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
			defer cancel()
			require.NoError(t, withHTTP.Stop(ctx, slog.Default()))

			_, err := net.DialTimeout("tcp", withHTTP.address, 50*time.Millisecond)
			require.Error(t, err)
			require.Contains(t, err.Error(), "connection refused")
		})
	})

	t.Run("all services get mounted", func(t *testing.T) {
		withHTTP := NewWithHTTP("localhost:0")

		helloCalled := false
		byeCalled := false
		service1 := func(handle func(path string, handler func(http.ResponseWriter, *http.Request))) {
			handle("/hello", func(w http.ResponseWriter, r *http.Request) {
				helloCalled = true
			})

			handle("/bye", func(w http.ResponseWriter, r *http.Request) {
				byeCalled = true
			})
		}

		foobarCalled := false
		service2 := func(handle func(path string, handler func(http.ResponseWriter, *http.Request))) {
			handle("/foobar", func(w http.ResponseWriter, r *http.Request) {
				foobarCalled = true
				require.NoError(t, r.Body.Close())
			})
		}

		go func() {
			require.NoError(t, withHTTP.Start(context.Background(), slog.Default(), service1, service2))
		}()
		defer func() {
			require.NoError(t, withHTTP.Stop(context.Background(), slog.Default()))
		}()

		select {
		case <-withHTTP.StartedChan():
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "timed out waiting for server to start")
		}

		_, err := http.Get(fmt.Sprintf("http://%s/hello", withHTTP.address))
		require.NoError(t, err)
		require.True(t, helloCalled, "hello should have been called")

		_, err = http.Get(fmt.Sprintf("http://%s/bye", withHTTP.address))
		require.NoError(t, err)
		require.True(t, byeCalled, "bye should have been called")

		_, err = http.Get(fmt.Sprintf("http://%s/foobar", withHTTP.address))
		require.NoError(t, err)
		require.True(t, foobarCalled, "foobar should have been called")
	})
}
