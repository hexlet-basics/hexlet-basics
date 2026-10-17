-- ADR-0019, state transitions slice: publishing, unpublishing and archiving a
-- blog post become operations of their own (blog_posts.publish,
-- blog_posts.unpublish, blog_posts.archive), no longer a field of the update.
-- The state used to be settable through the update, so every role holding
-- blog_posts.update gets the three new Permissions: nobody gains or loses access.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "grants"."role_id", "new"."permission", now(), now()
FROM "staff_member_role_grants" AS "grants"
CROSS JOIN (VALUES ('blog_posts.publish'), ('blog_posts.unpublish'), ('blog_posts.archive')) AS "new" ("permission")
WHERE "grants"."permission" = 'blog_posts.update'
ON CONFLICT ("role_id", "permission") DO NOTHING;
