#!/bin/sh
# Every internal/ directory is a singular package name that no ent/<x> package
# shares (a file needing both would need an import alias). A name is taken as
# plural when it ends in "s" but not "ss" (so mass nouns such as "progress"
# pass). `events`, `jobs` and `handlers` stay plural on purpose: the singular
# would shadow the commonest local variable names.
#
# scripts/package_names_exclusions.txt lists today's offenders; each rename
# ticket removes its entries. An entry that no longer offends fails too, so the
# list can only shrink.
set -eu

repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
internal_dir=${INTERNAL_DIR:-$repo_root/internal}
ent_dir=${ENT_DIR:-$repo_root/ent}
exclusions_file=${EXCLUSIONS_FILE:-$repo_root/scripts/package_names_exclusions.txt}

allowed="events jobs handlers"
failed=0

contains() { # contains <list> <word>
  for w in $1; do [ "$w" = "$2" ] && return 0; done
  return 1
}

excluded=$(grep -v '^[[:space:]]*\(#\|$\)' "$exclusions_file" || true)
offenders=""

for path in "$internal_dir"/*/; do
  name=$(basename "$path")
  offends=""
  case "$name" in
    *ss) ;;
    *s) contains "$allowed" "$name" || offends="plural" ;;
  esac
  if [ -d "$ent_dir/$name" ]; then offends="collides with ent/$name"; fi
  [ -n "$offends" ] || continue
  offenders="$offenders $name"
  if ! contains "$excluded" "$name"; then
    echo "internal/$name: $offends; use a singular name that matches no ent package (or add it to scripts/package_names_exclusions.txt only for an unmigrated package)" >&2
    failed=1
  fi
done

for name in $excluded; do
  if ! contains "$offenders" "$name"; then
    echo "scripts/package_names_exclusions.txt lists $name, which no longer offends: remove the entry" >&2
    failed=1
  fi
done

exit "$failed"
