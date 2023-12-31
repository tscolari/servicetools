package logging

import (
	"log/slog"
	"os"
)

// Default returns the default logger.
var Default = func() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
}
