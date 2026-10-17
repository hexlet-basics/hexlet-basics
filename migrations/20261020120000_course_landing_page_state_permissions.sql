-- ADR-0018, state transitions slice: publishing, unpublishing and archiving a
-- landing page become operations of their own (course_landing_pages.publish,
-- course_landing_pages.unpublish, course_landing_pages.archive), no longer a field of the update.
-- The state used to be settable through the update, so every role holding
-- course_landing_pages.update gets the three new Permissions: nobody gains or loses access.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "grants"."role_id", "new"."permission", now(), now()
FROM "staff_member_role_grants" AS "grants"
CROSS JOIN (VALUES ('course_landing_pages.publish'), ('course_landing_pages.unpublish'), ('course_landing_pages.archive')) AS "new" ("permission")
WHERE "grants"."permission" = 'course_landing_pages.update'
ON CONFLICT ("role_id", "permission") DO NOTHING;
