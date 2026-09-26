-- A user likes a blog post once. Legacy enforced that per browser session
-- (`session[:blog_post_likes]`), which the Go stack does not have, and the table
-- carried the rule nowhere: two non-unique indexes, on `blog_post_id` and on
-- `user_id`. The same user in two sessions could like a post twice, so
-- production may hold duplicate (blog_post_id, user_id) rows. The like handler
-- inserts conflict-tolerantly, which needs the constraint to be real.
--
-- Additive only (ADR-0015 rollback window): this migration declares the
-- invariant and does not clean up. If duplicates exist, creating the index
-- fails and the deploy stops — count them first and collapse them deliberately
-- as a cutover runbook step (#812). Likes have no dependents, so collapsing is
-- deleting all but the lowest id per pair. Rows with a NULL user_id (legacy
-- guests, if any) stay distinct under PostgreSQL's NULL semantics.
CREATE UNIQUE INDEX IF NOT EXISTS "index_blog_post_likes_on_blog_post_id_and_user_id"
ON "blog_post_likes" ("blog_post_id", "user_id");
