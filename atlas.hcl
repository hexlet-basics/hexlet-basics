# Atlas owns the database schema going forward. The first migration
# (migrations/…_baseline.sql) is the legacy Rails schema captured once via
# `atlas migrate diff` — Rails is retired, so all further changes are added with
# `atlas migrate new <name>` (hand-authored SQL) and applied by `make test-migrate`.
#
# `dev` is a throwaway Postgres that Atlas spins up to plan/lint migrations; it is
# never the target DB. `url` (the DB to apply to) is overridden per environment
# with `--url` — see the Makefile.

env "local" {
  url = "postgres://postgres:postgres@127.0.0.1:54330/code_basics_test?sslmode=disable&search_path=public"
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }
}
