package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
)

// NewWithHTTP returns a WithHTTP object configured with the given address.
func NewWithHTTP(address string) *WithHTTP {
	return &WithHTTP{
		address:     address,
		mutex:       new(sync.Mutex),
		startedChan: make(chan struct{}),
	}
}

// WithHTTP adds an HTTP server capability to another struct.
// It allows endpoints to be "registed" during Start, and provides
// a Stop method for shutting down the server.
// Once WithHTTP is ready to listen, it will send a signal to the
// channel returned by the StartedChan method.
type WithHTTP struct {
	address     string
	started     bool
	startedChan chan struct{}

	mutex  *sync.Mutex
	server *http.Server
	mux    *http.ServeMux
}

// HTTPRegisterFunc defines the functions that can be passed to Start
// in order to register new endpoints in the internal mux.
type HTTPRegisterFunc func(handle func(path string, handler func(http.ResponseWriter, *http.Request)))

// ConfigureHTTP is the hook used by the cmd package to inject the
// WithHTTP object in the host struct. This must be implemented by the host struct.
func (s *WithHTTP) ConfigureHTTP(*WithHTTP) {
	panic("ConfigureHTTP must be implemented")
}

// Start will register all given registerFuncs to the internal mux, bind
// the internal HTTP server to the listening address and block until the server shuts down.
// To wait for the server to start, the channel in the StartedChan() method can be used.
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

	s.address = listener.Addr().String()
	s.mux = http.NewServeMux()
	for _, registerFunc := range registerFuncs {
		// inject "interceptors" here wrapping mix.Handle
		registerFunc(s.mux.HandleFunc)
	}

	s.server = &http.Server{
		Handler: s.mux,
	}

	s.started = true

	logger.Info("starting HTTP Server", "address", listener.Addr().String())

	s.mutex.Unlock()
	close(s.startedChan)

	if err := s.server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("http server returned an error: %w", err)
		}
	}

	return nil
}

// Stop will gracefully stop the internal HTTP server.
// This will cause the Start function to return.
func (s *WithHTTP) Stop(ctx context.Context, logger *slog.Logger) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.server.Shutdown(ctx)
}

// StartedChan returns a channel that can be used to observe if
// the server has started or not.
func (s *WithHTTP) StartedChan() <-chan struct{} {
	return s.startedChan
}
