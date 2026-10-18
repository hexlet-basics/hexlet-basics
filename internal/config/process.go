package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

// LoadDotenv fills the environment from a local .env file, for a developer's
// checkout only: it does nothing unless APP_ENV is unset or "development", so
// a production process or a test never reads a dotenv file. godotenv never
// overrides a variable that is already set, so exported variables win. A
// missing file is not an error: .env is a convenience.
//
// It is a function of its own, not part of LoadServer/LoadWorker, so loading
// configuration stays a pure read of the environment; cmd/* calls it first.
func LoadDotenv() error {
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" && appEnv != string(Development) {
		return nil
	}
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return oops.Wrapf(err, "load .env")
	}
	return nil
}

// ColorOverride reports the NO_COLOR / FORCE_COLOR conventions: set says
// whether either is present, enabled what it asks for. NO_COLOR wins. It lives
// here because the environment is read only in this package.
func ColorOverride() (enabled, set bool) {
	if os.Getenv("NO_COLOR") != "" {
		return false, true
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true, true
	}
	return false, false
}
