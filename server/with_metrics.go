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

// NewWithMetrics returns a WithMetrics object configured with address.
func NewWithMetrics(address string) *WithMetrics {
	return &WithMetrics{
		address: address,
	}
}

// WithMetrics implements a simple HTTP server that responds to the `/metrics` endpoint
// with exposed prometheus metrics.
type WithMetrics struct {
	address string

	listener net.Listener
	server   *http.Server
}

// StartMetrics will start the HTTP metrics server and block
// until the StopMetrics method is called.
// This also takes an optional healthHandler, which must implement the healthcheck interface -
// When provided, that will allow this server to respond to liveness and readiness probes.
func (h *WithMetrics) StartMetrics(logger *slog.Logger, healthHandler healthcheck.Handler) error {

	lis, err := net.Listen("tcp", h.address)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	mux := http.NewServeMux()
	if healthHandler != nil {
		mux.Handle("/", healthHandler)
	}

	mux.Handle("/metrics", promhttp.Handler())

	h.server = &http.Server{Handler: mux}
	h.listener = lis

	logger.Info("starting Metrics Server", "address", h.listener.Addr().String())

	if err := h.server.Serve(h.listener); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("metrics server returned an error: %w", err)
		}
	}

	return nil
}

// StopMetrics will stop the Metrics server and cause Start() to unblock.
func (h *WithMetrics) StopMetrics(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

// ConfigureMetrics is the hook used by the cmd package to inject the
// WithMetrics object in the host struct. This must be implemented by the host struct.
func (h *WithMetrics) ConfigureMetrics(*WithMetrics) {
	panic("ConfigureMetrics must be implemented")
}
