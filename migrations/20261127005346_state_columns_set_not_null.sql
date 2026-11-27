-- Second half of 20261127005345_state_columns_not_null: validate the checks (a SHARE UPDATE
-- EXCLUSIVE lock, reads and writes go on), then SET NOT NULL, which PostgreSQL proves from
-- the validated helper check without a scan, and drop the helper, which NOT NULL now states.
-- The values checks stay.

SET lock_timeout = '5s';
SET statement_timeout = '60s';

ALTER TABLE "banners" VALIDATE CONSTRAINT "banners_state_values";

ALTER TABLE "blog_posts" VALIDATE CONSTRAINT "blog_posts_state_values";
ALTER TABLE "blog_posts" VALIDATE CONSTRAINT "blog_posts_state_not_null";
ALTER TABLE "blog_posts" ALTER COLUMN "state" SET NOT NULL;
-- The helper was only the scan-free path to the NOT NULL above.
-- squawk-ignore ban-drop-constraint
ALTER TABLE "blog_posts" DROP CONSTRAINT "blog_posts_state_not_null";

ALTER TABLE "reviews" VALIDATE CONSTRAINT "reviews_state_values";
ALTER TABLE "reviews" VALIDATE CONSTRAINT "reviews_state_not_null";
ALTER TABLE "reviews" ALTER COLUMN "state" SET NOT NULL;
-- The helper was only the scan-free path to the NOT NULL above.
-- squawk-ignore ban-drop-constraint
ALTER TABLE "reviews" DROP CONSTRAINT "reviews_state_not_null";

ALTER TABLE "language_landing_pages" VALIDATE CONSTRAINT "language_landing_pages_state_values";
ALTER TABLE "language_landing_pages" VALIDATE CONSTRAINT "language_landing_pages_state_not_null";
ALTER TABLE "language_landing_pages" ALTER COLUMN "state" SET NOT NULL;
-- The helper was only the scan-free path to the NOT NULL above.
-- squawk-ignore ban-drop-constraint
ALTER TABLE "language_landing_pages" DROP CONSTRAINT "language_landing_pages_state_not_null";

ALTER TABLE "languages" VALIDATE CONSTRAINT "languages_progress_values";
ALTER TABLE "languages" VALIDATE CONSTRAINT "languages_progress_not_null";
ALTER TABLE "languages" ALTER COLUMN "progress" SET NOT NULL;
-- The helper was only the scan-free path to the NOT NULL above.
-- squawk-ignore ban-drop-constraint
ALTER TABLE "languages" DROP CONSTRAINT "languages_progress_not_null";
