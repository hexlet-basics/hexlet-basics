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
	secret := Secret(secretValue)
	holder := struct {
		Key Secret `json:"key"`
	}{secret}

	var logs bytes.Buffer
	slog.New(slog.NewJSONHandler(&logs, nil)).Info("loaded", "secret", secret, "config", holder)
	asJSON, err := json.Marshal(holder)
	require.NoError(t, err)

	outputs := map[string]string{
		"%s":        fmt.Sprintf("%s", secret),
		"%v":        fmt.Sprintf("%v", secret),
		"%+v":       fmt.Sprintf("%+v", secret),
		"%#v":       fmt.Sprintf("%#v", secret),
		"%q":        fmt.Sprintf("%q", secret),
		"%x":        fmt.Sprintf("%x", secret),
		"struct %v": fmt.Sprintf("%v", holder),
		"Sprint":    fmt.Sprint(secret),
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
