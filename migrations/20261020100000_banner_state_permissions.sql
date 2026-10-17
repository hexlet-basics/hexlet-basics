-- ADR-0017, state transitions slice: publishing, unpublishing and archiving a
-- banner become operations of their own (banners.publish,
-- banners.unpublish, banners.archive), no longer a field of the update.
-- The state used to be settable through the update, so every role holding
-- banners.update gets the three new Permissions: nobody gains or loses access.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "grants"."role_id", "new"."permission", now(), now()
FROM "staff_member_role_grants" AS "grants"
CROSS JOIN (VALUES ('banners.publish'), ('banners.unpublish'), ('banners.archive')) AS "new" ("permission")
WHERE "grants"."permission" = 'banners.update'
ON CONFLICT ("role_id", "permission") DO NOTHING;
