package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
)

func NewWithHTTP(address string) *WithHTTP {
	return &WithHTTP{
		address:     address,
		mutex:       new(sync.Mutex),
		startedChan: make(chan struct{}),
	}
}

type WithHTTP struct {
	address     string
	started     bool
	startedChan chan struct{}

	mutex  *sync.Mutex
	server *http.Server
	mux    *http.ServeMux
}

type HTTPRegisterFunc func(handle func(path string, handler http.Handler))

func (s *WithHTTP) ConfigureHTTP(*WithHTTP) {
	panic("ConfigureHTTP must be implemented")
}

func (s *WithHTTP) Start(ctx context.Context, logger *slog.Logger, registerFuncs ...HTTPRegisterFunc) error {
	s.mutex.Lock()

	if s.started {
		s.mutex.Unlock()
		return fmt.Errorf("the server was already started")
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.mutex.Unlock()
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.mux = http.NewServeMux()
	for _, registerFunc := range registerFuncs {
		// inject "interceptors" here wrapping mix.Handle
		registerFunc(s.mux.Handle)
	}

	s.server = &http.Server{
		Handler: s.mux,
	}

	s.started = true

	logger.Info("starting HTTP Server", "address", listener.Addr().String())

	s.mutex.Unlock()
	s.startedChan <- struct{}{}

	if err := s.server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("http server returned an error: %w", err)
		}
	}

	return nil
}

func (s *WithHTTP) Stop(ctx context.Context, logger *slog.Logger) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.server.Shutdown(ctx)

	return nil
}

func (s *WithHTTP) StartedChan() <-chan struct{} {
	return s.startedChan
}
