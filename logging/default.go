package logging

import (
	"log/slog"
	"os"
)

// Default returns a logger at default configuration.
var Default = func() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
}
