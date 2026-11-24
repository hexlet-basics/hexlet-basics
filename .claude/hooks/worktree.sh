#!/usr/bin/env bash
# Claude Code WorktreeCreate/WorktreeRemove hook: keeps git worktrees in `<repo>/.worktrees/<name>`
# instead of the harness default `.claude/worktrees/<name>`, the same layout across our repositories.
#
#   worktree.sh create   # stdin: {"name": ..., "cwd": ...}; stdout: the worktree path, nothing else
#   worktree.sh remove   # stdin: {"worktree_path": ...}
#
# A new tree gets branch `worktree-<name>` (the harness default) from the remote default branch.
# Remove refuses a tree with uncommitted changes (it stays) and deletes the branch only when its
# work is already in the default branch. All git output goes to stderr.
set -euo pipefail

input=$(cat)
field() { printf '%s' "$input" | jq -r --arg k "$1" '.[$k] // empty'; }

# The primary checkout, even when the path is inside one of its worktrees.
primary() { dirname "$(git -C "$1" rev-parse --path-format=absolute --git-common-dir)"; }

default_ref() {
  git -C "$1" symbolic-ref -q --short refs/remotes/origin/HEAD \
    || { git -C "$1" rev-parse -q --verify origin/main >/dev/null && echo origin/main; } \
    || echo origin/master
}

case "${1:-}" in
  create)
    name=$(field name)
    [ -n "$name" ] || { echo "WorktreeCreate: no worktree name in the event" >&2; exit 1; }
    root=$(primary "$(field cwd)")
    dir="$root/.worktrees/$name"
    branch="worktree-$name"
    if [ ! -e "$dir/.git" ]; then
      git -C "$root" fetch -q origin >&2 || true
      if git -C "$root" show-ref -q --verify "refs/heads/$branch"; then
        git -C "$root" worktree add "$dir" "$branch" >&2
      else
        # --no-track: otherwise the branch tracks the default branch and a bare `git push` targets it.
        git -C "$root" worktree add --no-track -b "$branch" "$dir" "$(default_ref "$root")" >&2
      fi
    fi
    echo "$dir"
    ;;
  remove)
    path=$(field worktree_path)
    [ -n "$path" ] && [ -e "$path" ] || exit 0
    root=$(primary "$path")
    branch=$(git -C "$path" branch --show-current)
    git -C "$root" worktree remove "$path" >&2
    if [ -n "$branch" ] && git -C "$root" merge-base --is-ancestor "$branch" "$(default_ref "$root")" 2>/dev/null; then
      git -C "$root" branch -D "$branch" >&2
    fi
    ;;
  *)
    echo "usage: worktree.sh create|remove < hook-event.json" >&2
    exit 2
    ;;
esac
