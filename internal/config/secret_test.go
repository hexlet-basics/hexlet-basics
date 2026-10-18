package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

const secretValue = "hunter2-very-secret"

func TestSecretNeverPrintsItsValue(t *testing.T) {
	s := Secret(secretValue)
	holder := struct{ Key Secret }{s}

	var logs bytes.Buffer
	slog.New(slog.NewJSONHandler(&logs, nil)).Info("loaded", "secret", s, "config", holder)
	asJSON, err := json.Marshal(holder)
	require.NoError(t, err)

	outputs := map[string]string{
		"%s":        fmt.Sprintf("%s", s), //nolint:perfsprint // the verb is the case under test
		"%v":        fmt.Sprintf("%v", s),
		"%+v":       fmt.Sprintf("%+v", s),
		"%#v":       fmt.Sprintf("%#v", s),
		"%q":        fmt.Sprintf("%q", s),
		"%x":        fmt.Sprintf("%x", s),
		"struct %v": fmt.Sprintf("%v", holder),
		"Sprint":    fmt.Sprint(s),
		"slog":      logs.String(),
		"json":      string(asJSON),
	}
	for verb, out := range outputs {
		require.NotContains(t, out, secretValue, verb)
		require.Contains(t, out, "REDACTED", verb)
	}
}

func TestSecretRevealReturnsTheValue(t *testing.T) {
	require.Equal(t, secretValue, Secret(secretValue).Reveal())
}
