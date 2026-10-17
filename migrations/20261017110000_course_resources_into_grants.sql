-- ADR-0017, third slice: the course side of the admin (course categories,
-- courses, course lessons, course lesson reviews, lesson progress and course
-- landing pages) moves from the create/read/update/delete matrix to declared
-- Permissions. The matrix named these resources with the legacy Rails values
-- (language_categories, languages, language_lessons, language_lesson_reviews,
-- language_lesson_members, language_landing_pages); the Permissions use the
-- domain terms. Each bit becomes the grant that opens exactly the operations the
-- old policy let it open (can_index: list and get, can_create: create,
-- can_update: update and review, can_destroy: delete), so nobody gains or loses
-- access. Bits that opened nothing (courses had no delete, lessons only a list
-- and a review) are not carried.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "role_id", 'course_categories.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_categories' AND "can_index"
UNION ALL
SELECT "role_id", 'course_categories.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_categories' AND "can_create"
UNION ALL
SELECT "role_id", 'course_categories.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_categories' AND "can_update"
UNION ALL
SELECT "role_id", 'course_categories.delete', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_categories' AND "can_destroy"
UNION ALL
SELECT "role_id", 'courses.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'languages' AND "can_index"
UNION ALL
SELECT "role_id", 'courses.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'languages' AND "can_create"
UNION ALL
SELECT "role_id", 'courses.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'languages' AND "can_update"
UNION ALL
SELECT "role_id", 'course_lessons.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_lessons' AND "can_index"
UNION ALL
SELECT "role_id", 'course_lessons.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_lessons' AND "can_update"
UNION ALL
SELECT "role_id", 'course_lesson_reviews.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_lesson_reviews' AND "can_index"
UNION ALL
SELECT "role_id", 'lesson_progress.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_lesson_members' AND "can_index"
UNION ALL
SELECT "role_id", 'course_landing_pages.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_landing_pages' AND "can_index"
UNION ALL
SELECT "role_id", 'course_landing_pages.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_landing_pages' AND "can_create"
UNION ALL
SELECT "role_id", 'course_landing_pages.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_landing_pages' AND "can_update"
UNION ALL
SELECT "role_id", 'course_landing_pages.delete', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'language_landing_pages' AND "can_destroy";

-- The matrix no longer decides these resources; leaving the rows would show a
-- second, dead copy of the right.
DELETE FROM "staff_member_role_permissions" WHERE "resource" IN ('language_categories', 'languages', 'language_lessons', 'language_lesson_reviews', 'language_lesson_members', 'language_landing_pages');
