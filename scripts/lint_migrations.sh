#!/usr/bin/env bash
# A shipped migration never changes, and a new one is lint-clean for PostgreSQL
# locking and downtime (squawk, .squawk.toml). Compared with the merge base of
# MIGRATIONS_BASE (default origin/main), so a push to main that adds nothing passes.
#
# A change to an existing migration needs the marker in the commit message
# ("Migration-Edit: <why>") so the exception is a recorded human decision.
set -euo pipefail

base_ref=${MIGRATIONS_BASE:-origin/main}
if ! git rev-parse --verify --quiet "$base_ref" >/dev/null; then
  echo "lint:migrations: $base_ref is not fetched (use fetch-depth: 0 in CI)" >&2
  exit 1
fi
base=$(git merge-base "$base_ref" HEAD)

changed=$(git diff --name-status --no-renames "$base" -- 'migrations/*.sql' | awk '$1 != "A" { print $2 }')
if [ -n "$changed" ] && ! git log --format=%B "$base..HEAD" | grep -q '^Migration-Edit: '; then
  echo "an existing migration changed or was removed (shipped migrations are immutable):" >&2
  echo "$changed" >&2
  echo "add a 'Migration-Edit: <why>' line to the commit message if that is intended" >&2
  exit 1
fi

added=$(git diff --name-only --diff-filter=A --no-renames "$base" -- 'migrations/*.sql')
if [ -n "$added" ]; then
  # shellcheck disable=SC2086 # one path per line, none with spaces
  squawk $added
fi
