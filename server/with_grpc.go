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

func NewWithGRPC(address string) *WithGRPC {
	return &WithGRPC{
		address:     address,
		mutex:       new(sync.Mutex),
		startedChan: make(chan struct{}),
	}
}

type GRPCRegisterFunc func(*grpc.Server)

type WithGRPC struct {
	address     string
	started     bool
	startedChan chan struct{}

	mutex  *sync.Mutex
	server *grpc.Server
}

func (s *WithGRPC) ConfigureGRPC(*WithGRPC) {
	panic("ConfigureGRPC must be implemented")
}

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
		grpc.ChainUnaryInterceptor(
			grpcsrv.LoggerInterceptor(logger),
			grpcsrv.LoggerAnnotationInterceptor,
		),
	)

	for _, registerFunc := range registerFuncs {
		registerFunc(s.server)
	}

	s.started = true
	logger.Info("starting GRPC Server", "address", listener.Addr().String())
	s.mutex.Unlock()

	s.startedChan <- struct{}{}

	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("grpc server returned an error: %w", err)
	}

	return nil
}

func (s *WithGRPC) StartedChan() <-chan struct{} {
	return s.startedChan
}

func (s *WithGRPC) Stop(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.server.GracefulStop()

	return nil
}
