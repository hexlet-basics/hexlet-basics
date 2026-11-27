-- The Publication State of banners, blog posts, reviews and course landing pages, and the
-- Course Readiness of courses, become typed, required values that only the generated state
-- move writes (ADR-0029). The Rails schema had them as nullable strings (banners excepted:
-- NOT NULL DEFAULT 'draft'), and the code read a NULL as a draft. A move is now a
-- conditional update on the current value, which a NULL would never match, and the ent
-- schema types the column with the vocabulary, which the database now enforces too.
--
-- This file backfills the NULLs as drafts, sets the default, and adds two NOT VALID checks
-- per column: the values of the vocabulary (kept), and NOT NULL (a helper). The next file
-- validates them in its own transaction, so no table is scanned under an ACCESS EXCLUSIVE
-- lock, and turns the helper into the column's NOT NULL.

SET lock_timeout = '5s';
SET statement_timeout = '60s';

ALTER TABLE "banners" ADD CONSTRAINT "banners_state_values" CHECK ("state" IN ('draft', 'published', 'archived')) NOT VALID;

UPDATE "blog_posts" SET "state" = 'draft' WHERE "state" IS NULL;
ALTER TABLE "blog_posts" ALTER COLUMN "state" SET DEFAULT 'draft';
ALTER TABLE "blog_posts" ADD CONSTRAINT "blog_posts_state_not_null" CHECK ("state" IS NOT NULL) NOT VALID;
ALTER TABLE "blog_posts" ADD CONSTRAINT "blog_posts_state_values" CHECK ("state" IN ('draft', 'published', 'archived')) NOT VALID;

UPDATE "reviews" SET "state" = 'draft' WHERE "state" IS NULL;
ALTER TABLE "reviews" ALTER COLUMN "state" SET DEFAULT 'draft';
ALTER TABLE "reviews" ADD CONSTRAINT "reviews_state_not_null" CHECK ("state" IS NOT NULL) NOT VALID;
ALTER TABLE "reviews" ADD CONSTRAINT "reviews_state_values" CHECK ("state" IN ('draft', 'published', 'archived')) NOT VALID;

UPDATE "language_landing_pages" SET "state" = 'draft' WHERE "state" IS NULL;
ALTER TABLE "language_landing_pages" ALTER COLUMN "state" SET DEFAULT 'draft';
ALTER TABLE "language_landing_pages" ADD CONSTRAINT "language_landing_pages_state_not_null" CHECK ("state" IS NOT NULL) NOT VALID;
ALTER TABLE "language_landing_pages" ADD CONSTRAINT "language_landing_pages_state_values" CHECK ("state" IN ('draft', 'published', 'archived')) NOT VALID;

UPDATE "languages" SET "progress" = 'draft' WHERE "progress" IS NULL;
ALTER TABLE "languages" ALTER COLUMN "progress" SET DEFAULT 'draft';
ALTER TABLE "languages" ADD CONSTRAINT "languages_progress_not_null" CHECK ("progress" IS NOT NULL) NOT VALID;
ALTER TABLE "languages" ADD CONSTRAINT "languages_progress_values" CHECK ("progress" IN ('draft', 'in_development', 'completed')) NOT VALID;
