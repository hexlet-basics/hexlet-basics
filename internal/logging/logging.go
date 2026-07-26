// Package logging builds the application logger.
package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// New returns an slog.Logger backed by a tinted handler writing to stderr.
func New(level slog.Level) *slog.Logger {
	return slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{Level: level}))
}
