package server

import (
	context "context"
	slog "log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/testhelpers"
)

func Test_WithWorker(t *testing.T) {
	t.Run("started channel: multiple listeners are allowed", func(t *testing.T) {
		listener1Ok := false
		listener2Ok := false

		withWorker := NewWithWorker()

		go func() {
			<-withWorker.StartedChan()
			listener1Ok = true
		}()

		go func() {
			<-withWorker.StartedChan()
			listener2Ok = true
		}()

		testhelpers.Constantly(t, func() bool {
			return (!listener1Ok) && (!listener2Ok)
		}, 300*time.Millisecond, 100*time.Millisecond, "channel listener changed")

		go func() {
			require.NoError(t, withWorker.Start(context.Background(), slog.Default()))
		}()

		require.Eventually(t, func() bool {
			return listener1Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 1 didn't return")

		require.Eventually(t, func() bool {
			return listener2Ok
		}, 500*time.Millisecond, 50*time.Millisecond, "channel listener 2 didn't return")

	})

	t.Run("all tasks are mounted and running", func(t *testing.T) {
		withWorker := NewWithWorker()

		taskOneChan := make(chan struct{}, 100)
		taskOne := func(ctx context.Context, logger *slog.Logger) error {
			logger.Info("task one started")
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Millisecond):
					taskOneChan <- struct{}{}
				}
			}
		}

		taskTwoChan := make(chan struct{}, 100)
		taskTwo := func(ctx context.Context, logger *slog.Logger) error {
			logger.Info("task two started")
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Millisecond):
					taskTwoChan <- struct{}{}
				}
			}
		}

		go func() {
			require.NoError(t, withWorker.Start(context.Background(), slog.Default(), taskOne, taskTwo))
		}()

		select {
		case <-withWorker.StartedChan():
		case <-time.After(100 * time.Millisecond):
			require.Fail(t, "timed out waiting for worker to start")
		}

		require.Eventually(t, func() bool {
			return len(taskOneChan) > 2
		}, 50*time.Millisecond, 5*time.Millisecond)

		require.Eventually(t, func() bool {
			return len(taskTwoChan) > 2
		}, 50*time.Millisecond, 5*time.Millisecond)

		withWorker.Stop(context.Background(), slog.Default())

		taskOneChanLenght := len(taskOneChan)
		taskTwoChanLenght := len(taskTwoChan)

		testhelpers.Constantly(t, func() bool {
			return taskOneChanLenght == len(taskOneChan)
		}, 50*time.Millisecond, 5*time.Millisecond)

		testhelpers.Constantly(t, func() bool {
			return taskTwoChanLenght == len(taskTwoChan)
		}, 50*time.Millisecond, 5*time.Millisecond)

	})
}
