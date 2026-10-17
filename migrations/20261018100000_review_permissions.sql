-- ADR-0018, delegation slice: reviewing a course and reviewing a course lesson
-- become Permissions of their own (courses.review, course_lessons.review),
-- separate from updating. Until now the review was opened by the update
-- Permission, so every role holding it gets the matching review Permission and
-- nobody gains or loses access.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "role_id", 'courses.review', now(), now() FROM "staff_member_role_grants"
  WHERE "permission" = 'courses.update'
UNION ALL
SELECT "role_id", 'course_lessons.review', now(), now() FROM "staff_member_role_grants"
  WHERE "permission" = 'course_lessons.update'
ON CONFLICT ("role_id", "permission") DO NOTHING;

-- The review was the only operation course_lessons.update opened, so the grant
-- now names nothing and the catalog no longer lists it. It is replaced by the
-- review grant above rather than left behind as a dead right.
DELETE FROM "staff_member_role_grants" WHERE "permission" = 'course_lessons.update';
