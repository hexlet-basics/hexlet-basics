-- An Enrollment is a learner's record in one Course, and never created twice
-- for the same pair (CONTEXT.md). `language_members` carried that invariant
-- nowhere: two non-unique indexes, on `user_id` and on `language_id`, and a
-- find-or-create in the legacy app that read before it wrote. Two concurrent
-- requests could each insert. Everything in the progress epic inserts
-- conflict-tolerantly, which needs the constraint to be real in PostgreSQL.
--
-- FIXME(#765): the production duplicate count has NOT been read yet. This
-- migration collapses in a single transaction, which is right for a handful of
-- duplicates and wrong for thousands — at that volume the collapse belongs in
-- its own ticket with batched updates and its own verification pass. Run
-- scripts/oneoff/765-duplicate-enrollments/count.sql against production (or a
-- replica) and confirm before this reaches a production database. Query 4 of
-- that script returning anything but 0 also invalidates the re-point below.
--
-- The collapse keeps the OLDEST Enrollment per (user, course), ordered by
-- `created_at, id`. The `id` tie-break is not decoration: a racy find-or-create
-- is exactly what produces identical `created_at` values, and without it the
-- winner would differ between a dry run and the real one.
--
-- Re-pointing cannot collide: `language_lesson_members` is unique on
-- (`user_id`, `lesson_id`), which does not mention `language_member_id` at all,
-- so moving a loser's rows onto the winner changes only the FK. No
-- deduplication of Lesson Progress rows is needed.
--
-- Every statement is idempotent against a database with no duplicates: the CTE
-- selects nothing, the UPDATE and DELETE touch nothing, and the index creation
-- is guarded.

-- Re-point the losing Enrollments' Lesson Progress rows onto the winner.
WITH ranked AS (
  SELECT
    id,
    first_value(id) OVER (
      PARTITION BY user_id, language_id
      ORDER BY created_at, id
    ) AS winner_id
  FROM language_members
)
UPDATE language_lesson_members llm
SET language_member_id = r.winner_id
FROM ranked r
WHERE llm.language_member_id = r.id
  AND r.id <> r.winner_id;

-- Sum the losers' denormalized finished-lesson counters onto the winner. The
-- counter is not the source of Completion (ADR-0012 computes that from the
-- current Version), but it stays consistent with the rows it claims to count.
WITH ranked AS (
  SELECT
    id,
    finished_lessons_count,
    first_value(id) OVER (
      PARTITION BY user_id, language_id
      ORDER BY created_at, id
    ) AS winner_id
  FROM language_members
),
losers AS (
  SELECT winner_id, sum(finished_lessons_count) AS extra
  FROM ranked
  WHERE id <> winner_id
  GROUP BY winner_id
)
UPDATE language_members lm
SET finished_lessons_count = lm.finished_lessons_count + losers.extra
FROM losers
WHERE lm.id = losers.winner_id;

-- Delete the losers, now that nothing references them.
WITH ranked AS (
  SELECT
    id,
    first_value(id) OVER (
      PARTITION BY user_id, language_id
      ORDER BY created_at, id
    ) AS winner_id
  FROM language_members
)
DELETE FROM language_members
WHERE id IN (SELECT id FROM ranked WHERE id <> winner_id);

CREATE UNIQUE INDEX IF NOT EXISTS "index_language_members_on_user_id_and_language_id"
ON "language_members" ("user_id", "language_id");
