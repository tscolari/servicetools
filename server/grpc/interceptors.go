package grpc

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"

	"github.com/gogo/status"
	"github.com/tscolari/servicetools/logging"
)

const (
	loggerFieldMethod     = "grpc_method"
	loggerFieldStatusCode = "status_code"
)

// LoggerInterceptor adds a logger to the context and annotates it.
func LoggerInterceptor(logger *slog.Logger) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = logging.ToContext(ctx, logger)
		return handler(ctx, req)
	}
}

// LoggerAnnotationInterceptor adds annotations to the logger based on the request.
// It also logs debugging information for every request.
func LoggerAnnotationInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	logger := logging.FromContext(ctx)
	logger = logger.With(loggerFieldMethod, info.FullMethod)

	now := time.Now()
	logger.Debug("request started")

	ctx = logging.ToContext(ctx, logger)
	h, err := handler(ctx, req)

	errCode := status.Code(err).String()
	logger.Debug(
		"request finished",
		"error", err,
		loggerFieldStatusCode, errCode,
		"duration", time.Since(now).String(),
	)

	return h, err
}
