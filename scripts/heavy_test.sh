#!/bin/sh
# Behaviour test for heavy: serialization, re-entrancy, timeout, exit status, CI bypass.
set -eu

heavy="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)/heavy"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
cd "$work"
git init -q .
unset CI HEAVY_LOCK_HELD HEAVY_LOCK_TIMEOUT

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

# Two holders must not overlap: each records start and end, then we check order.
sh "$heavy" sh -c 'echo a-start >>log; sleep 1; echo a-end >>log' &
sleep 0.3
sh "$heavy" sh -c 'echo b-start >>log; echo b-end >>log' 2>err &
wait
[ "$(tr '\n' ' ' <log)" = "a-start a-end b-start b-end " ] || fail "runs overlapped: $(tr '\n' ' ' <log)"
grep -q 'waiting for the heavy-task lock' err || fail "no waiting notice"

# The command's exit status passes through.
status=0
sh "$heavy" sh -c 'exit 3' || status=$?
[ "$status" -eq 3 ] || fail "exit status $status, wanted 3"

# Re-entrant: a nested heavy does not deadlock.
sh "$heavy" sh "$heavy" true || fail "nested heavy failed"

# Timeout: a held lock makes a waiter fail with a clear message.
sh "$heavy" sleep 3 &
sleep 0.3
status=0
HEAVY_LOCK_TIMEOUT=1 sh "$heavy" true 2>err || status=$?
wait
[ "$status" -eq 75 ] || fail "timeout status $status, wanted 75"
grep -q 'gave up' err || fail "no timeout message"

# CI: no locking at all.
sh "$heavy" sleep 3 &
sleep 0.3
CI=1 sh "$heavy" true || fail "CI run failed"
wait
echo "heavy_test: ok"
