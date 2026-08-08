-- Count duplicate course enrollments — see issue #765 (parent #755).
--
-- FIXME(#765): the numbers below still have to be read off PRODUCTION. Nobody has
-- run this against prod yet, so the duplicate count, the worst-case multiplicity
-- and the affected participation count are all still unknown. Until they are, the
-- migration that adds `UNIQUE (user_id, language_id)` to `language_members` cannot
-- be written: it is not known whether the duplicate collapse fits inside that
-- migration or needs a ticket of its own. Record the output on #765 and fill in
-- the TODOs in the report template at the bottom of this file.
--
-- Why this exists: `language_members` (a course enrollment) has only two
-- non-unique indexes, on `user_id` and on `language_id`. There is no unique index
-- on the pair, and the legacy find-or-create was racy, so two concurrent requests
-- could each insert an enrollment for the same learner and course.
--
--   course enrollment    = language_members        (user_id, language_id)
--   lesson participation = language_lesson_members (language_member_id -> language_members.id)
--
-- READ-ONLY: no writes and no locks beyond shared read locks. Queries 1-3 are
-- full scans with a group-by over `language_members`, so prefer a read replica.
--
--   psql "$PROD_DATABASE_URL" -f scripts/count_duplicate_enrollments.sql
--
-- Verified by scripts/count_duplicate_enrollments_verify.sql, which shadows both
-- tables with a fixture whose answers are known by hand.

\pset pager off

-- 1. Headline numbers, as one labeled result set so the whole block can be pasted
--    back without correlating anything by hand.
--
--    rows_to_delete is the number that actually sizes the migration, and it is
--    NOT the same as duplicate_pairs: one pair holding four enrollments
--    contributes one to duplicate_pairs but three to rows_to_delete.
WITH pair_counts AS (
  SELECT user_id, language_id, count(*) AS c
  FROM language_members
  GROUP BY user_id, language_id
)
SELECT 'duplicate_pairs'    AS metric, count(*)::text                AS value FROM pair_counts WHERE c > 1
UNION ALL
SELECT 'worst_multiplicity',        coalesce(max(c), 0)::text                 FROM pair_counts
UNION ALL
SELECT 'rows_to_delete',            coalesce(sum(c - 1), 0)::text             FROM pair_counts
UNION ALL
SELECT 'total_enrollments',         count(*)::text                            FROM language_members;

-- 2. The full multiplicity distribution, not just the maximum. This is what
--    answers "a handful or thousands", and it distinguishes a lone pathological
--    pair from a systematic pattern — a max of 9 means something very different
--    if exactly one pair has it than if a thousand pairs do.
SELECT c AS enrollments_per_pair, count(*) AS pairs
FROM (
  SELECT count(*) AS c
  FROM language_members
  GROUP BY user_id, language_id
) t
GROUP BY c
ORDER BY c;

-- 3. Lesson participations hanging off the LOSING enrollments — the rows a
--    collapse would have to re-point onto the winner.
--
--    The winner is the oldest enrollment, ordered by `created_at, id`. The `id`
--    tie-break is required rather than decorative: the racy find-or-create is
--    exactly the thing that produces identical `created_at` values, and with ties
--    present `row_number()` would pick a different winner between runs, so a
--    dry run would not predict what the real migration did.
--
--    Cross-check: `losing_enrollments` here and `rows_to_delete` from query 1 are
--    computed by different means — a window function versus a sum over group
--    counts — and must come back equal. If production returns two different
--    numbers, the assumption behind these queries is wrong, not production.
WITH ranked AS (
  SELECT id,
         row_number() OVER (PARTITION BY user_id, language_id
                            ORDER BY created_at, id) AS rn
  FROM language_members
),
losers AS (SELECT id FROM ranked WHERE rn > 1)
SELECT
  (SELECT count(*) FROM losers)                                     AS losing_enrollments,
  (SELECT count(*) FROM language_lesson_members llm
     JOIN losers l ON l.id = llm.language_member_id)                 AS participations_to_repoint;

-- 4. Consistency check, and the one result that can force the answer on its own.
--
--    `language_lesson_members` denormalizes `user_id` and `language_id` alongside
--    its `language_member_id` FK. Every row should agree with the enrollment it
--    points at. If any row does not, re-pointing participations stops being a
--    mechanical `UPDATE` and the collapse needs its own verification pass no
--    matter how few duplicates there are.
--
--    Expected: 0.
SELECT count(*) AS participations_disagreeing_with_their_enrollment
FROM language_lesson_members llm
JOIN language_members lm ON lm.id = llm.language_member_id
WHERE llm.user_id <> lm.user_id OR llm.language_id <> lm.language_id;

--------------------------------------------------------------------------------
-- Report template for #765. Fill in the TODOs from the output above.
--------------------------------------------------------------------------------
--
--   duplicate (learner, course) pairs .................. TODO
--   worst-case multiplicity ........................... TODO
--   losing enrollment rows to delete ................... TODO
--   participations on losing enrollments ............... TODO
--   total enrollments (denominator) .................... TODO
--   participations disagreeing with their enrollment .... TODO
--   multiplicity distribution .......................... TODO
--
-- Why volume is the only open question. `language_lesson_members` carries
-- `UNIQUE (user_id, lesson_id)` — the unique index does not mention
-- `language_member_id` at all. Re-pointing a loser's participations onto the
-- winner changes only `language_member_id` and leaves `user_id` and `lesson_id`
-- alone, so it cannot collide with the winner's own participations; two rows for
-- one `(user_id, lesson_id)` cannot exist on either side of the collapse. That
-- makes the collapse three statements with no conflict handling and no
-- participation deduplication:
--
--   1. pick the winner per (user_id, language_id) by `created_at, id`;
--   2. UPDATE language_lesson_members SET language_member_id = <winner>
--        WHERE language_member_id IN (<losers>);
--   3. sum finished_lessons_count onto the winner, then delete the losers.
--
-- The decision rule, fixed here BEFORE the numbers are read so that the verdict
-- is not a post-hoc reading of whatever came back:
--
--   * any nonzero result from query 4  -> its own ticket, whatever the volume,
--     because the re-point is no longer mechanical;
--   * otherwise rows_to_delete <= ~100 -> collapse inside the index migration.
--     At that size the collapse is three statements in the same transaction that
--     adds the unique index, and creating the index IS the verification: it
--     cannot succeed while a duplicate remains;
--   * otherwise                        -> its own ticket, with batched updates,
--     a counter reconciliation check, and its own verification pass before the
--     index is added.
--
--   verdict ............................................ TODO
