package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

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

type WithWorker struct {
	ctx       context.Context
	cancelCtx func()

	mutex       *sync.Mutex
	started     bool
	startedChan chan error
	wg          *sync.WaitGroup
}

func (w *WithWorker) ConfigureWorker(*WithWorker) {
	panic("ConfigureWorker must be implemented")
}

func (w *WithWorker) Start(ctx context.Context, logger *slog.Logger, tasks ...WorkerTaskFunc) error {
	w.mutex.Lock()

	if w.started {
		w.mutex.Unlock()
		return fmt.Errorf("the worker was already started")
	}

	w.started = true
	w.wg = new(sync.WaitGroup)

	for _, task := range tasks {
		w.wg.Add(1)
		task := task

		go func() {
			defer w.wg.Done()
			if err := task(ctx, logger); err != nil {
				w.startedChan <- err
			}
		}()
	}

	logger.Info("starting Worker server")
	w.mutex.Unlock()
	w.startedChan <- nil
	w.wg.Wait()

	return nil
}

func (w *WithWorker) Stop(ctx context.Context, logger *slog.Logger) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.cancelCtx()
	return nil
}

func (w *WithWorker) StartedChan() <-chan error {
	return w.startedChan
}
