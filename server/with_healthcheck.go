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

// NewWithHealthcheck returns a WithHealthcheck object configured with address and metrics.
func NewWithHealthcheck(address string, metrics bool) *WithHealthcheck {
	return &WithHealthcheck{
		address: address,
		metrics: metrics,
	}
}

// WithHealthcheck implements a simple HTTP server that has `/live` and `/ready` endpoints.
// If started with `metrics: true` it will also expose metrics through the `/metrics` endpoint.
type WithHealthcheck struct {
	address string
	metrics bool

	listener net.Listener
	server   *http.Server
}

// StartHealthcheck will start the HTTP healthcheck server and block
// until the StopHealthcheck method is called.
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

// StopHealthcheck will stop the Healthcheck server and cause Start() to unblock.
func (h *WithHealthcheck) StopHealthcheck(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

// ConfigureHealthcheck is the hook used by the cmd package to inject the
// WithHealthcheck object. Services using WithHealthcheck must overwrite this method.
func (h *WithHealthcheck) ConfigureHealthcheck(*WithHealthcheck) {
	panic("ConfigureHealthcheck must be implemented")
}
