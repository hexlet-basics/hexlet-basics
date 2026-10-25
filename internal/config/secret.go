package config

import (
	"fmt"
	"log/slog"
)

// redacted is what a Secret prints in place of its value.
const redacted = "[REDACTED]"

// Secret is a setting whose value must not reach a log line, an error message
// or a JSON dump by accident: a configuration struct is exactly the kind of
// value that gets logged whole while debugging a start-up failure. Every
// printing path (fmt verbs, slog, JSON, text) yields a placeholder; the value
// is available only through Reveal, so each use is visible in the code and can
// be searched for. It is a string underneath so the env decoder and the
// validator tags treat it like any other string setting.
type Secret string

// Reveal returns the value. Call it only where the secret is handed to the
// component that needs it (a signer, an API client), never to print it.
func (s Secret) Reveal() string { return string(s) }

// String implements fmt.Stringer, which covers %s and %v.
func (Secret) String() string { return redacted }

// GoString covers %#v.
func (Secret) GoString() string { return redacted }

// Format covers every verb (%q, %x, %d...), which would otherwise fall back to
// the underlying string.
// A fmt.State write error has nowhere to go (fmt itself drops it), so it is
// written with fmt.Fprint rather than the State's Write.
func (Secret) Format(f fmt.State, _ rune) { fmt.Fprint(f, redacted) }

// LogValue implements slog.LogValuer, so a Secret in a structured log field is
// redacted too.
func (Secret) LogValue() slog.Value { return slog.StringValue(redacted) }

// MarshalJSON redacts the value in JSON.
func (Secret) MarshalJSON() ([]byte, error) { return []byte(`"` + redacted + `"`), nil }

// MarshalText redacts the value in every text encoding (YAML, TOML, text/template).
func (Secret) MarshalText() ([]byte, error) { return []byte(redacted), nil }
