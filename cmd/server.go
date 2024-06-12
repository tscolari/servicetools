package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/tscolari/servicetools/database"
	"github.com/tscolari/servicetools/logging"
	"github.com/tscolari/servicetools/server"
)

// Server defines the interface that this command uses to start/stop.
type Server interface {
	Start(context.Context, *slog.Logger) error
	Stop(context.Context, *slog.Logger) error
}

// HasGRPC means the Server has gRPC capability.
type HasGRPC interface {
	ConfigureGRPC(*server.WithGRPC)
}

// HasHTTP means the Server has HTTP capability.
type HasHTTP interface {
	ConfigureHTTP(*server.WithHTTP)
}

// HasWorker means the Server has worker capability.
type HasWorker interface {
	ConfigureWorker(*server.WithWorker)
}

// HasMetrics means the Server has metrics capability.
type HasMetrics interface {
	ConfigureMetrics(*server.WithMetrics)
}

// HasDatabase means the Server has database capability.
type HasDatabase interface {
	ConfigureDatabase(*server.WithDB)
}

// HasReaderDatabase means the Server has database reader capability.
type HasReaderDatabase interface {
	ConfigureReaderDatabase(*server.WithRDB)
}

// CanServer injects the "server" (or start) subcommand to another command.
// It will start the given Server based on the capabilities that it implements.
func CanServer(rootCmd *cobra.Command, srv Server) {
	serverToRun = srv

	// Enable only the flags that the given server supports:

	if _, ok := serverToRun.(HasDatabase); ok {
		serverCmd.PersistentFlags().StringVar(&serverDBEnvPrefix, "db-env-prefix", "DATABASE", "prefix to env variables with DB configuration")
	}

	if _, ok := serverToRun.(HasReaderDatabase); ok {
		serverCmd.PersistentFlags().StringVar(&serverRDBEnvPrefix, "reader-db-env-prefix", "DATABASE_READER", "prefix to env variables with READER DB configuration")
	}

	if _, ok := serverToRun.(HasGRPC); ok {
		serverCmd.PersistentFlags().StringVar(&serverGRPCAddress, "grpc-address", "localhost:0", "listening address for GRPC connections")
	}

	if _, ok := serverToRun.(HasHTTP); ok {
		serverCmd.PersistentFlags().StringVar(&serverHTTPAddress, "http-address", "localhost:0", "listening address for HTTP connections")
	}

	if _, ok := serverToRun.(HasMetrics); ok {
		serverCmd.PersistentFlags().StringVar(&serverMetricsAddress, "metrics-address", "localhost:0", "listening address for metrics")
	}

	rootCmd.AddCommand(serverCmd)
}

var (
	serverToRun Server

	serverGRPCAddress    string
	serverHTTPAddress    string
	serverMetricsAddress string
	serverDBEnvPrefix    string
	serverRDBEnvPrefix   string
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"start"},
	Short:   "Starts the server",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		logger := logging.Default()

		if grpcSrv, ok := serverToRun.(HasGRPC); ok {
			withGRPC := server.NewWithGRPC(serverGRPCAddress)
			grpcSrv.ConfigureGRPC(withGRPC)
		}

		if httpSrv, ok := serverToRun.(HasHTTP); ok {
			withHTTP := server.NewWithHTTP(serverHTTPAddress)
			httpSrv.ConfigureHTTP(withHTTP)
		}

		if workerSrv, ok := serverToRun.(HasWorker); ok {
			withWorker := server.NewWithWorker()
			workerSrv.ConfigureWorker(withWorker)
		}

		if metricsSrv, ok := serverToRun.(HasMetrics); ok {
			withMetrics := server.NewWithMetrics(serverMetricsAddress)
			metricsSrv.ConfigureMetrics(withMetrics)
		}

		if withDBSrv, ok := serverToRun.(HasDatabase); ok {
			dbConfig, err := database.ConfigFromEnv(serverDBEnvPrefix)
			if err != nil {
				logger.Error("failed to load database configuration", "error", err)
				return fmt.Errorf("failed to generate DB configuration: %w", err)
			}

			withDB, err := server.NewWithDB(dbConfig)
			if err != nil {
				logger.Error("failed to configure database", "error", err)
				return fmt.Errorf("failed to configure DB: %w", err)
			}
			withDBSrv.ConfigureDatabase(withDB)
		}

		if withRDBSrv, ok := serverToRun.(HasReaderDatabase); ok {
			dbConfig, err := database.ConfigFromEnv(serverRDBEnvPrefix)
			if err != nil {
				logger.Error("failed to load database configuration", "error", err)
				return fmt.Errorf("failed to generate DB configuration: %w", err)
			}

			withRDB, err := server.NewWithRDB(dbConfig)
			if err != nil {
				logger.Error("failed to configure reader database", "error", err)
				return fmt.Errorf("failed to configure Reader DB: %w", err)
			}
			withRDBSrv.ConfigureReaderDatabase(withRDB)
		}

		go func() {
			stopSignal := make(chan os.Signal, 1)
			signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)

			select {
			case s := <-stopSignal:
				logger.Info("signal received, exiting", "signal", s.String())
			case <-ctx.Done():
				logger.Info("server context was closed, exiting")
			}

			if err := serverToRun.Stop(ctx, logger); err != nil {
				logger.Error("attempt to stop server failed", "error", err)
			}
		}()

		if err := serverToRun.Start(ctx, logger); err != nil {
			logger.Error("server failed", "error", err)
			return err
		}

		return nil
	},
}
