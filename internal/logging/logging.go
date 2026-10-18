// Package logging builds the application logger.
package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"hexletbasics/internal/config"
)

// New returns an slog.Logger backed by a tinted handler writing to stderr.
func New(level slog.Level) *slog.Logger {
	return slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
		Level:   level,
		NoColor: !colorEnabled(os.Stderr),
	}))
}

// colorEnabled keeps ANSI escapes out of anything that is not a terminal: tint
// colours by default, which turns container logs and log files into noise.
// NO_COLOR and FORCE_COLOR follow the usual conventions, the latter for dev
// daemons whose output goes to a file a human still reads in a terminal.
func colorEnabled(f *os.File) bool {
	if enabled, set := config.ColorOverride(); set {
		return enabled
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
