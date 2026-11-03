-- schemalens found two kinds of index debt in the schema inherited from the Rails app.
--
-- 1. A foreign key with no index on its referencing column: deleting or updating a referenced
--    row scans the whole referencing table. Seven columns get an index.
-- 2. An index whose columns are the leading columns of another index of the table, which
--    serves the same lookups and costs a write each: seven are dropped, each named with the
--    index that covers it.

CREATE INDEX "index_language_module_descriptions_on_language_id" ON "language_module_descriptions" ("language_id");
CREATE INDEX "index_language_module_version_infos_on_version_id" ON "language_module_version_infos" ("version_id");
CREATE INDEX "index_oauth_access_tokens_on_client_id" ON "oauth_access_tokens" ("client_id");
CREATE INDEX "index_oauth_authorization_requests_on_client_id" ON "oauth_authorization_requests" ("client_id");
CREATE INDEX "index_oauth_authorization_requests_on_user_id" ON "oauth_authorization_requests" ("user_id");
CREATE INDEX "index_oauth_refresh_tokens_on_client_id" ON "oauth_refresh_tokens" ("client_id");
CREATE INDEX "index_user_accounts_on_user_id" ON "user_accounts" ("user_id");

-- Covered by index_blog_post_likes_on_blog_post_id_and_user_id.
DROP INDEX IF EXISTS "index_blog_post_likes_on_blog_post_id";
-- Covered by index_language_members_on_user_id_and_language_id.
DROP INDEX IF EXISTS "index_language_members_on_user_id";
-- Covered by language_modules_language_id_slug_uniq.
DROP INDEX IF EXISTS "language_modules_language_id_index";
-- Covered by taggings_idx (tag_id first).
DROP INDEX IF EXISTS "index_taggings_on_tag_id";
-- Covered by taggings_idy and taggings_taggable_context_idx (taggable_id first).
DROP INDEX IF EXISTS "index_taggings_on_taggable_id";
-- Covered by index_taggings_on_taggable_type_and_taggable_id.
DROP INDEX IF EXISTS "index_taggings_on_taggable_type";
-- Covered by index_taggings_on_tagger_id_and_tagger_type.
DROP INDEX IF EXISTS "index_taggings_on_tagger_id";
