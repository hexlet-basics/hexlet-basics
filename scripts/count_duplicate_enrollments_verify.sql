-- Verification for scripts/count_duplicate_enrollments.sql — see issue #765.
--
-- Running the counting queries against an empty or duplicate-free database proves
-- nothing: every result is zero whether the SQL is right or wrong. This harness
-- shadows both tables with temp tables holding a fixture small enough to check by
-- hand, so the expected output is known in advance. Temp tables take search_path
-- precedence over the real ones, so the counting script picks these up unchanged.
--
-- Run both files in one session, this one first:
--
--   psql "$DATABASE_URL" --single-transaction \
--     -f scripts/count_duplicate_enrollments_verify.sql \
--     -f scripts/count_duplicate_enrollments.sql
--
-- Fixture — (user, course) -> enrollment ids:
--
--   (1, 10) -> 100            single enrollment; must NOT count as a duplicate
--   (1, 20) -> 200, 201       multiplicity 2, winner 200 (older created_at)
--   (2, 10) -> 300, 301, 302  multiplicity 3, winner 300
--   (3, 10) -> 400, 401       multiplicity 2 with IDENTICAL created_at, so the
--                             winner is 400 only because of the `id` tie-break.
--                             Drop `id` from the ORDER BY and this pair's winner
--                             becomes arbitrary — this row is why it is there.
--
-- Expected output:
--   duplicate_pairs     3        (the three multi-enrollment pairs)
--   worst_multiplicity  3        (pair (2, 10))
--   rows_to_delete      4        (1 + 2 + 1 losers)
--   total_enrollments   8
--   distribution        1 -> 1 pair, 2 -> 2 pairs, 3 -> 1 pair
--   losing_enrollments  4        (201, 301, 302, 401)
--   participations_to_repoint     3   (on 201, 302, 401 — NOT the 2 on winner 200)
--   participations_disagreeing    1   (row 6, planted below)

CREATE TEMP TABLE language_members (
  id bigint PRIMARY KEY,
  created_at timestamptz NOT NULL,
  finished_lessons_count integer NOT NULL DEFAULT 0,
  language_id bigint NOT NULL,
  state varchar,
  updated_at timestamptz NOT NULL,
  user_id bigint NOT NULL
);

CREATE TEMP TABLE language_lesson_members (
  id bigint PRIMARY KEY,
  created_at timestamptz NOT NULL,
  language_id bigint NOT NULL,
  language_member_id bigint NOT NULL,
  lesson_id bigint NOT NULL,
  messages_count integer DEFAULT 0,
  state varchar,
  updated_at timestamptz NOT NULL,
  user_id bigint NOT NULL
);

INSERT INTO language_members (id, created_at, updated_at, user_id, language_id) VALUES
  (100, '2024-01-01', '2024-01-01', 1, 10),
  (200, '2024-01-01', '2024-01-01', 1, 20),
  (201, '2024-02-01', '2024-02-01', 1, 20),
  (300, '2024-01-01', '2024-01-01', 2, 10),
  (301, '2024-02-01', '2024-02-01', 2, 10),
  (302, '2024-03-01', '2024-03-01', 2, 10),
  -- identical created_at: only the id tie-break makes the winner deterministic
  (400, '2024-05-05 12:00:00+00', '2024-05-05', 3, 10),
  (401, '2024-05-05 12:00:00+00', '2024-05-05', 3, 10);

INSERT INTO language_lesson_members
  (id, created_at, updated_at, language_id, language_member_id, lesson_id, user_id) VALUES
  -- on the winner 200: must NOT be counted as needing a re-point
  (1, '2024-01-01', '2024-01-01', 20, 200, 1, 1),
  (2, '2024-01-01', '2024-01-01', 20, 200, 2, 1),
  -- on losers 201, 302, 401: must be counted
  (3, '2024-01-01', '2024-01-01', 20, 201, 3, 1),
  (4, '2024-01-01', '2024-01-01', 10, 302, 4, 2),
  (5, '2024-01-01', '2024-01-01', 10, 401, 5, 3),
  -- disagrees with enrollment 100, which belongs to user 1: query 4 must catch it
  (6, '2024-01-01', '2024-01-01', 10, 100, 6, 99);
