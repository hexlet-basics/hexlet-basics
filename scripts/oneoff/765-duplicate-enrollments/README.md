# Count duplicate Enrollments (#765)

**One-off. Delete this whole directory once #765 is closed.** It is deliberately
not wired into the `Makefile` — nothing here belongs to the permanent toolchain,
and `scripts/oneoff/` exists so a throwaway measurement does not accumulate
targets that outlive it.

## Why

`language_members` (the Enrollment table) has only two non-unique indexes, on
`user_id` and on `language_id`. There is no unique index on the pair, and the
legacy find-or-create was racy, so two concurrent requests could each insert an
Enrollment for the same learner and Course.

The migration in #755 that adds `UNIQUE (user_id, language_id)` cannot be written
until the size of the existing mess is known: a handful of duplicates collapse
inside that migration, thousands need their own ticket with its own verification
pass. These queries produce that number.

## Running it

```sh
# The real thing — production, or a read replica for preference.
psql "$PROD_DATABASE_URL" -f count.sql

# Prove the queries are correct first (see below).
psql "$DATABASE_URL" --single-transaction -f fixture.sql -f count.sql
```

`count.sql` is **read-only**: no writes, and no locks beyond shared reads. Queries
1–3 full-scan `language_members` with a group-by, which is why a replica is
preferable.

## Why there is a fixture

On a duplicate-free database every counting query returns zero whether the SQL is
right or wrong, so running it somewhere clean proves nothing. `fixture.sql` shadows
both tables with temp tables holding input whose answers are computable by hand,
then prints those expected answers immediately above the real output.

Temp tables take `search_path` precedence, so `count.sql` runs **unchanged** and is
itself the thing under test — no query is restated, which is what stops the check
from drifting away from the script that actually runs against production.

The fixture discriminates rather than merely being non-zero: one pair has an
identical `created_at`, so `lesson_progress_to_repoint = 3` only holds if the
`id` tie-break elects the right winner. Drop `id` from the `ORDER BY` and the
expected answer changes.

## Vocabulary

`CONTEXT.md` names the concept **Enrollment**, and the per-Lesson record
**Lesson Progress**; Course Membership and Lesson Membership are on their
`_Avoid_` lines. These files said Membership until #766 flipped the glossary the
other way round — the words changed, none of the SQL did, because table and
column names (`language_members`, `language_lesson_members`,
`language_member_id`) keep the legacy vocabulary by design.

## What to do with the output

`count.sql` ends with a report template and a decision rule that was fixed
*before* any number was read, so the verdict follows from the output instead of
being argued after the fact. Fill in the `TODO`s, post them on #765, and record the
recommendation the rule yields.
