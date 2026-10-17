-- ADR-0019, state transitions slice: marking a course as draft, in development
-- or completed become operations of their own (courses.mark_draft,
-- courses.mark_in_development, courses.mark_completed), no longer a field of
-- the update. The readiness used to be settable through the update, so every
-- role holding courses.update gets the three new Permissions: nobody gains or
-- loses access.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "grants"."role_id", "new"."permission", now(), now()
FROM "staff_member_role_grants" AS "grants"
CROSS JOIN (VALUES ('courses.mark_draft'), ('courses.mark_in_development'), ('courses.mark_completed')) AS "new" ("permission")
WHERE "grants"."permission" = 'courses.update'
ON CONFLICT ("role_id", "permission") DO NOTHING;
