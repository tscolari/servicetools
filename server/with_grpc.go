package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/grpc"

	grpcsrv "github.com/tscolari/servicetools/server/grpc"
)

// NewWithGRPC returns a WithGRPC object set to listen at the given address.
func NewWithGRPC(address string, options ...grpc.ServerOption) *WithGRPC {
	return &WithGRPC{
		address:     address,
		mutex:       new(sync.Mutex),
		startedChan: make(chan struct{}),
		options:     options,
	}
}

// GRPCRegisterFunc is used as arguments to the Start method.
// It exposes the internal gRPC server and allow gRPC services to register to it.
type GRPCRegisterFunc func(*grpc.Server)

// WithGRPC defines the gRPC server capability.
type WithGRPC struct {
	address     string
	started     bool
	startedChan chan struct{}

	options []grpc.ServerOption
	mutex   *sync.Mutex
	server  *grpc.Server
}

// Start will bind the internal gRPC server to the address and execute all
// given registerFuncs.
// This will block until the server is stopped (using Stop()).
func (s *WithGRPC) Start(ctx context.Context, logger *slog.Logger, registerFuncs ...GRPCRegisterFunc) error {
	if s.started {
		return nil
	}

	s.mutex.Lock()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.mutex.Unlock()
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.server = grpc.NewServer(
		append(s.options,
			grpc.ChainUnaryInterceptor(
				grpcsrv.LoggerInterceptor(logger),
				grpcsrv.LoggerAnnotationInterceptor,
			),
		)...,
	)

	for _, registerFunc := range registerFuncs {
		registerFunc(s.server)
	}

	s.address = listener.Addr().String()
	s.started = true
	logger.Info("starting GRPC Server", "address", s.address)
	s.mutex.Unlock()

	close(s.startedChan)

	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("grpc server returned an error: %w", err)
	}

	return nil
}

// StartedChan can be used by a caller to block until the server has started.
// Once the server has started, the channel will be closed and unblocked.
func (s *WithGRPC) StartedChan() <-chan struct{} {
	return s.startedChan
}

// Stop will gracefully stop the internal gRPC Server.
func (s *WithGRPC) Stop(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.server.GracefulStop()
	return nil
}

// ConfigureGRPC is the hook used by the cmd package to inject the
// WithGRPC object in the host struct. This must be implemented by the host struct.
func (s *WithGRPC) ConfigureGRPC(*WithGRPC) {
	panic("ConfigureGRPC must be implemented")
}
