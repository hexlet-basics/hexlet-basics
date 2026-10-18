#!/bin/sh
# Behaviour test for lint_package_names.sh: runs it over throw-away trees.
set -eu

script="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)/lint_package_names.sh"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM

# run <internal dirs> <ent dirs> <exclusions>: exit status of the check.
run() {
  rm -rf "$work/internal" "$work/ent"
  mkdir -p "$work/internal" "$work/ent"
  for d in $1; do mkdir "$work/internal/$d"; done
  for d in $2; do mkdir "$work/ent/$d"; done
  echo "$3" | tr ' ' '\n' >"$work/exclusions"
  INTERNAL_DIR="$work/internal" ENT_DIR="$work/ent" EXCLUSIONS_FILE="$work/exclusions" sh "$script" >"$work/out" 2>&1
}

expect() { # expect <ok|fail> <description> <command...>
  want=$1
  desc=$2
  shift 2
  if "$@"; then got=ok; else got=fail; fi
  if [ "$got" != "$want" ]; then
    echo "FAIL: $desc (wanted $want)" >&2
    cat "$work/out" >&2
    exit 1
  fi
}

expect ok "singular, allow-listed and mass nouns pass" run "account events jobs handlers progress dashboard catalog" "user" ""
expect fail "a new plural fails" run "books" "user" ""
expect fail "a name equal to an ent package fails" run "lead" "lead" ""
expect ok "an excluded offender passes" run "books" "user" "books"
expect fail "a stale exclusion fails" run "book" "user" "books"
