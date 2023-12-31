package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewWithHealthcheck(address string, metrics bool) *WithHealthcheck {
	return &WithHealthcheck{
		address: address,
		metrics: metrics,
	}
}

type WithHealthcheck struct {
	address string
	metrics bool

	listener net.Listener
	server   *http.Server
}

func (h *WithHealthcheck) ConfigureHealthcheck(*WithHealthcheck) {
	panic("ConfigureHealthcheck must be implemented")
}

func (h *WithHealthcheck) StartHealthcheck(logger *slog.Logger) error {

	lis, err := net.Listen("tcp", h.address)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	health := healthcheck.NewHandler()

	mux := http.NewServeMux()
	mux.Handle("/", health)

	if h.metrics {
		mux.Handle("/metrics", promhttp.Handler())
	}

	h.server = &http.Server{Handler: mux}
	h.listener = lis

	logger.Info("starting Healthcheck Server", "address", h.listener.Addr().String())

	if err := h.server.Serve(h.listener); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("healthcheck server returned an error: %w", err)
		}
	}

	return nil
}

func (h *WithHealthcheck) StopHealthcheck(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
