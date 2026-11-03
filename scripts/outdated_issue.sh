#!/usr/bin/env bash
# Keeps one open issue that lists the mise tools whose newest release is not the one mise.toml
# pins (Dependabot has no mise ecosystem), and closes it once every pin is current. Run weekly by
# .github/workflows/tools.yml with GH_TOKEN set; needs mise, gh and jq on the PATH.
set -euo pipefail

title='Update the mise tools'
label='ready-for-agent'

# --bump: an exact pin counts as satisfied without it, so a pinned tool would never show as behind.
outdated=$(mise outdated --local --bump --json)
rows=$(jq -r 'to_entries | sort_by(.key)[] | .value
  | "| \(.name) | \(.current) | " + (if .release_url then "[\(.latest)](\(.release_url))" else .latest end) + " |"' <<<"$outdated")

# The search matches on words, not exactly, so filter on the exact title.
number=$(gh issue list --state open --label "$label" --search "\"$title\" in:title" --json number,title \
  --jq "map(select(.title == \"$title\")) | first | .number // empty")

if [[ -z "$rows" ]]; then
  if [[ -n "$number" ]]; then
    gh issue close "$number" --comment 'mise.toml pins the newest release of every tool.'
  fi
  exit 0
fi

body="These tools are pinned in \`mise.toml\` and nothing moves the pins but a commit.

| Tool | Pinned | Newest |
|---|---|---|
$rows

Bump them with \`mise run upgrade\` (it rewrites \`mise.toml\`), read the release notes of a tool
that fails \`mise run lint\`, and land it with a \`build(deps):\` commit. Mind the pins that carry a
reason in \`mise.toml\` (atlas, hk with hk.pkl). The newest release may be withheld by mise's
release-age cutoff; if \`mise upgrade\` skips a tool, leave it for the next run."

if [[ -z "$number" ]]; then
  gh issue create --title "$title" --label "$label" --body "$body"
else
  gh issue edit "$number" --body "$body"
fi
