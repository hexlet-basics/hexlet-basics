// Package logging builds the application logger (ADR-0025, extending ADR-0007).
package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"

	"hexletbasics/internal/config"
)

// redacted replaces the value of an attribute whose key names a secret.
const redacted = "[redacted]"

// sensitiveKeyParts are lowercase fragments of an attribute key that mark its
// value as one that must never reach a log: credentials and personal data.
// Identity in a log line is user_id, never an email.
var sensitiveKeyParts = []string{
	"password", "passwd", "secret", "token", "authorization", "cookie", "email", "api_key", "apikey",
}

// New returns the logger of a process, writing to stderr at the given level.
func New(level slog.Level) *slog.Logger {
	return NewTo(os.Stderr, level)
}

// NewTo is New with the destination chosen: tint text when w is a terminal,
// JSON otherwise, so containers and log files get one machine-readable object
// per line and a developer gets readable output. Every line passes through the
// trace-correlation wrapper and the key-name redaction safety net.
func NewTo(w io.Writer, level slog.Level) *slog.Logger {
	var inner slog.Handler
	if f, ok := w.(*os.File); ok && readable(f) {
		inner = tint.NewTextHandler(w, &tint.Options{Level: level, NoColor: !colorEnabled(f), ReplaceAttr: redact})
	} else {
		inner = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level, ReplaceAttr: redact})
	}
	return slog.New(traceHandler{inner})
}

// redact is the safety net for a log call that passes a secret by mistake. The
// review standard is that none is passed; this catches the slip by key name.
func redact(_ []string, attr slog.Attr) slog.Attr {
	key := strings.ToLower(attr.Key)
	for _, part := range sensitiveKeyParts {
		if strings.Contains(key, part) {
			return slog.String(attr.Key, redacted)
		}
	}
	return attr
}

// traceHandler adds trace_id and span_id of the context's span to each line.
// It tests SpanContext().IsValid(), not IsRecording(): production samples 1%
// of traces (ADR-0007), and a recording check would leave the other 99% of log
// lines without the key that finds them. trace_id is the only correlation key.
type traceHandler struct{ slog.Handler }

func (h traceHandler) Handle(ctx context.Context, record slog.Record) error {
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		record.AddAttrs(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, record)
}

func (h traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return traceHandler{h.Handler.WithAttrs(attrs)}
}

func (h traceHandler) WithGroup(name string) slog.Handler {
	return traceHandler{h.Handler.WithGroup(name)}
}

// readable reports whether a person reads f: a terminal, or a file FORCE_COLOR
// says a human still reads (dev daemons). Everything else gets JSON.
func readable(f *os.File) bool {
	if forced, _ := config.ColorOverride(); forced {
		return true
	}
	return isTerminal(f)
}

// colorEnabled keeps ANSI escapes out of anything that is not a terminal: tint
// colours by default, which turns container logs and log files into noise.
// NO_COLOR and FORCE_COLOR follow the usual conventions, the latter for dev
// daemons whose output goes to a file a human still reads in a terminal.
func colorEnabled(f *os.File) bool {
	if enabled, set := config.ColorOverride(); set {
		return enabled
	}
	return isTerminal(f)
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
