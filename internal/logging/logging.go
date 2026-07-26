// Package logging builds the application logger.
package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// New returns an slog.Logger backed by a tinted handler writing to stderr.
func New(level slog.Level) *slog.Logger {
	handler := tint.NewHandler(os.Stderr, &tint.Options{Level: level})
	return slog.New(handler)
}
