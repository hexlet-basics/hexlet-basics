# Git workflow: main checkout stays on main, real work goes in a worktree

**Never switch branches in the main checkout.** It stays on `main`. Several
agent sessions and the owner share it, and its working tree usually holds
uncommitted changes that are not yours; a `git checkout` there breaks them or
the other sessions.

## Trivial work: commit on main

A change of a few lines in a file or two (a typo, a doc tweak, a one-line fix)
is committed directly on `main` in the main checkout.

## Anything bigger: a worktree

If the work goes beyond trivial (several files, a feature, a refactor, anything
you would want to test or iterate on), do it in its own worktree:

- Create it with the `EnterWorktree` tool, or
  `git worktree add .claude/worktrees/<topic> -b <topic> origin/main`.
  Worktrees live in `.claude/worktrees/` (gitignored), never in a sibling
  directory next to the main checkout.
- Work, run the checks and commit there, on its branch.
- Land it on `main`: `git fetch`, `git rebase origin/main`, then
  `git push origin HEAD:main`. No pull requests for agent work.
- Clean up: remove the worktree (`git worktree remove`) and delete the branch
  locally and on `origin` (`git branch -d`, `git push origin --delete`).

## Both cases

- **Stage only your own files, by path.** Never `git add -A` / `git add .` /
  `git commit -a`, and never stash, reset or restore files you did not change.
- **Rebase onto `origin/main` before pushing.** Another session may have pushed
  since you started. If your commit is already on `origin/main`, there is
  nothing to push.
- **`.git/index.lock` exists:** another session's git is running. Wait and
  retry; do not delete the lock while a git process is alive
  (`ps aux | grep '[g]it '`).
- **Push only when asked** ("заливай", "push"). Commits follow Conventional
  Commits (see `AGENTS.md`); hooks must pass, never `--no-verify`.
