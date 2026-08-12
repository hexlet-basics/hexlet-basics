-- Collapse duplicate Enrollments — issue #765, prerequisite of #757.
--
-- One-off. Delete it with the rest of this directory once #765 is closed.
-- `language_members` gained UNIQUE (user_id, language_id) in migration
-- 20260812013103. That migration deliberately does NOT clean up: creating the
-- index fails loudly if duplicates exist, and cleaning them is a decision with
-- its own volume question, not a side effect of a deploy.
--
-- Run count.sql first. This script is right for a handful of duplicates — it
-- collapses everything in one transaction, holding row locks for its duration.
-- For thousands, batch it instead. A non-zero `mismatched_lesson_progress` from
-- count.sql query 4 invalidates the re-point below entirely: the denormalized
-- user_id/language_id on a Lesson Progress row would disagree with the
-- Enrollment it is being moved to.
--
--   psql "$PROD_DATABASE_URL" --single-transaction -f collapse.sql
--
-- The winner is the OLDEST Enrollment per (user, course), ordered by
-- `created_at, id`. The id tie-break is not decoration: a racy find-or-create is
-- exactly what produces identical `created_at` values, and without it the winner
-- would differ between a dry run and the real one.
--
-- Re-pointing cannot collide: `language_lesson_members` is unique on
-- (`user_id`, `lesson_id`), which does not mention `language_member_id`, so
-- moving a loser's rows onto the winner changes only the FK. No deduplication of
-- Lesson Progress rows is needed.
--
-- Every statement is idempotent: run against a database with no duplicates, the
-- CTEs select nothing and nothing is written. Verified in
-- internal/store/enrollment_collapse_test.go, which executes THIS file.

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
