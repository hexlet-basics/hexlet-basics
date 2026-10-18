#!/bin/sh
# The Atlas Community version is pinned twice: `atlas-community` in mise.toml
# (the CLI developers and CI run) and the `arigaio/atlas:<v>-community` image in
# the Dockerfile (the CLI that applies migrations in production). A bump that
# moves only one lets the migration tool differ between environments, so this
# fails when they disagree (ADR-0014). Paths are arguments so the check can be
# run against other files; the defaults are the repository's own.

set -eu

repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
mise_file=${1:-$repo_root/mise.toml}
docker_file=${2:-$repo_root/Dockerfile}

tool=$(sed -n 's/^atlas-community[[:space:]]*=[[:space:]]*"\([^"]*\)".*/\1/p' "$mise_file")
image=$(sed -n 's/^FROM[[:space:]]*arigaio\/atlas:\([^[:space:]]*\)-community.*/\1/p' "$docker_file")

if [ -z "$tool" ]; then
  echo "lint:atlas-pin: no atlas-community pin found in $mise_file" >&2
  exit 1
fi
if [ -z "$image" ]; then
  echo "lint:atlas-pin: no arigaio/atlas:<version>-community image found in $docker_file" >&2
  exit 1
fi
if [ "$tool" != "$image" ]; then
  echo "lint:atlas-pin: mise.toml pins atlas-community $tool but the Dockerfile copies $image-community; move them together (ADR-0014)" >&2
  exit 1
fi
