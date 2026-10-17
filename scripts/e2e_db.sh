#!/usr/bin/env bash
# Recreates the end-to-end database: a fresh code_basics_e2e_test, migrated with
# atlas and loaded with the committed fixtures/, the same starting data the Go
# tests get. The name carries "test" because the fixture loader refuses to
# wipe a database without it. Reached through E2E_DB_HOST / E2E_DB_PORT (the
# local db daemon on 54330 by default, a service container in CI); the
# development database is never touched.
set -euo pipefail

host=${E2E_DB_HOST:-127.0.0.1}
port=${E2E_DB_PORT:-54330}
name=${E2E_DB_NAME:-code_basics_e2e_test}
export PGPASSWORD=${E2E_DB_PASSWORD:-postgres}

nc -z "$host" "$port" 2>/dev/null || { echo "e2e db: no PostgreSQL on $host:$port (mise run db:up, or set E2E_DB_HOST/E2E_DB_PORT)" >&2; exit 1; }

psql=psql
command -v psql >/dev/null || psql="mise x postgres@18 -- psql"
run() { $psql -X -q -h "$host" -p "$port" -U postgres "$@"; }

run -d postgres -c "drop database if exists \"$name\" with (force)" -c "create database \"$name\""

dsn="postgres://postgres:${PGPASSWORD}@${host}:${port}/${name}?sslmode=disable"
atlas migrate apply --env local --url "${dsn}&search_path=public" >/dev/null
go run ./cmd/fixtures -conn "$dsn" -dir fixtures

# The fixtures carry no password digests (the Go tests mint sessions directly),
# but a browser has to sign in. This is the bcrypt digest of the public,
# throw-away password "e2e-password", given to the fixture admin, to a staff
# editor and to a user who is neither.
# shellcheck disable=SC2016 # the dollars are part of the bcrypt digest, not expansions
digest='$2a$10$igu.GskpHZdOZgUlOQ7.5u/TZus5zank/0OQiVpwIUJkg4bJjL5Ge'
run -d "$name" -c "update users set password_digest = '$digest' where email in ('alice@example.com', 'bob@example.com', 'carol@example.com')"
echo "e2e db ready: $name on $host:$port"
