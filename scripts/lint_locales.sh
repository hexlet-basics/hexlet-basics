#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
catalog_dir="$repo_root/internal/localization/locales"
work_dir=$(mktemp -d)
trap 'rm -rf "$work_dir"' EXIT HUP INT TERM

for locale in en ru es; do
  catalog="$catalog_dir/active.$locale.json"
  if [ ! -f "$catalog" ]; then
    echo "backend locale catalog is missing: active.$locale.json" >&2
    exit 1
  fi
  cp "$catalog" "$work_dir"
done

cd "$repo_root"
go tool goi18n extract \
  -sourceLanguage en \
  -format json \
  -outdir "$work_dir" \
  internal/localization

if ! cmp -s "$catalog_dir/active.en.json" "$work_dir/active.en.json"; then
  echo "backend source catalog is stale; run 'make gen-locales'" >&2
  diff -u "$catalog_dir/active.en.json" "$work_dir/active.en.json" || true
  exit 1
fi

go tool goi18n merge \
  -sourceLanguage en \
  -format json \
  -outdir "$work_dir" \
  "$work_dir"/active.*.json

set -- "$work_dir"/translate.*.json
if [ -e "$1" ]; then
  echo "backend locale catalogs have missing or empty translations:" >&2
  for catalog in "$@"; do
    echo "  $(basename "$catalog")" >&2
  done
  exit 1
fi
