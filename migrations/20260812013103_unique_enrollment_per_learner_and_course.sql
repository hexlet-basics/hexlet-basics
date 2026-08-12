-- An Enrollment is a learner's record in one Course, and never created twice
-- for the same pair (CONTEXT.md). `language_members` carried that invariant
-- nowhere: two non-unique indexes, on `user_id` and on `language_id`, and a
-- find-or-create in the legacy app that read before it wrote. Two concurrent
-- requests could each insert. Everything in the progress epic inserts
-- conflict-tolerantly, which needs the constraint to be real in PostgreSQL.
--
-- This migration only declares the invariant; it does not clean up. If
-- duplicates exist, creating the index fails and the deploy stops — which is
-- the intended behaviour: collapsing rows is a decision with its own volume
-- question (a handful merge in one transaction, thousands need batching), not a
-- side effect of a deploy nobody is watching. The measurement and the collapse
-- live in scripts/oneoff/765-duplicate-enrollments/ (#765), to be run first and
-- deliberately.
CREATE UNIQUE INDEX IF NOT EXISTS "index_language_members_on_user_id_and_language_id"
ON "language_members" ("user_id", "language_id");
