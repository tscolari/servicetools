package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// NewWithWorker returns a new worker object.
func NewWithWorker() *WithWorker {
	return &WithWorker{
		mutex:       new(sync.Mutex),
		startedChan: make(chan error),
	}
}

// WorkerTaskFunc defines a function that starts a task.
// Tasks will be started in parallel, and they should exit when
// the context is canceled.
type WorkerTaskFunc func(ctx context.Context, logger *slog.Logger) error

// WithWorker implements simple worker capabilities.
// It can be started with a list of taks (WorkerTaskFunc), where each
// will be spawn in a goroutine.
type WithWorker struct {
	ctx       context.Context
	cancelCtx func()

	mutex       *sync.Mutex
	started     bool
	startedChan chan error
	wg          *sync.WaitGroup
}

// ConfigureWorker is the hook used by the cmd package to inject the
// WithWorker object in the host struct. This must be implemented by the host struct.
func (w *WithWorker) ConfigureWorker(*WithWorker) {
	panic("ConfigureWorker must be implemented")
}

// Start will start all the given tasks, and block until all them are finished.
// Once all tasks are started/scheduled, the channel from StartedChan() will unblock.
// Tasks are WorkerTaskFunc functions, and they should exit once the given context
// is canceled.
func (w *WithWorker) Start(ctx context.Context, logger *slog.Logger, tasks ...WorkerTaskFunc) error {
	w.mutex.Lock()

	if w.started {
		w.mutex.Unlock()
		return fmt.Errorf("the worker was already started")
	}

	w.started = true
	w.wg = new(sync.WaitGroup)

	taskCtx, cancel := context.WithCancel(ctx)
	w.cancelCtx = cancel

	for _, task := range tasks {
		w.wg.Add(1)
		task := task

		go func() {
			defer w.wg.Done()
			if err := task(taskCtx, logger); err != nil {
				w.startedChan <- err
			}
		}()
	}

	logger.Info("starting Worker server")
	w.mutex.Unlock()
	close(w.startedChan)
	w.wg.Wait()

	return nil
}

// Stop will signal to all internal tasks to stop, by canceling their internal contexts.
func (w *WithWorker) Stop(ctx context.Context, logger *slog.Logger) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.cancelCtx()
	return nil
}

// StartedChan returns a channel that can be used to inspect if all the tasks have
// been started.
func (w *WithWorker) StartedChan() <-chan error {
	return w.startedChan
}
