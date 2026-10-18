package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// dotenvDir makes the working directory one that holds a .env file.
func dotenvDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("DOTENV_PROBE=from-file\n"), 0o600))
	t.Chdir(dir)
	t.Setenv("DOTENV_PROBE", "")
	require.NoError(t, os.Unsetenv("DOTENV_PROBE"))
	// godotenv sets the variable with os.Setenv, which t.Setenv does not undo.
	t.Cleanup(func() { _ = os.Unsetenv("DOTENV_PROBE") })
}

func TestLoadDotenvReadsTheFileInDevelopment(t *testing.T) {
	for _, appEnv := range []string{"", "development"} {
		dotenvDir(t)
		t.Setenv("APP_ENV", appEnv)

		require.NoError(t, LoadDotenv())

		require.Equal(t, "from-file", os.Getenv("DOTENV_PROBE"), "APP_ENV=%q", appEnv)
		require.NoError(t, os.Unsetenv("DOTENV_PROBE"))
	}
}

func TestLoadDotenvIgnoresTheFileInProductionAndTest(t *testing.T) {
	for _, appEnv := range []string{"production", "test"} {
		dotenvDir(t)
		t.Setenv("APP_ENV", appEnv)

		require.NoError(t, LoadDotenv())

		require.Empty(t, os.Getenv("DOTENV_PROBE"), "APP_ENV=%s", appEnv)
	}
}
