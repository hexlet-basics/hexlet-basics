-- ADR-0017, second slice: banners, reviews, leads and messages move from the
-- create/read/update/delete matrix to declared Permissions. Each bit becomes the
-- grant that opens exactly the operations the old policy let that bit open
-- (can_index: list and get, can_create: create, can_update: update, can_destroy:
-- delete), so nobody gains or loses access. Leads and messages only ever had a
-- list, so only their read bit meant anything.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "role_id", 'banners.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'banners' AND "can_index"
UNION ALL
SELECT "role_id", 'banners.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'banners' AND "can_create"
UNION ALL
SELECT "role_id", 'banners.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'banners' AND "can_update"
UNION ALL
SELECT "role_id", 'banners.delete', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'banners' AND "can_destroy"
UNION ALL
SELECT "role_id", 'reviews.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'reviews' AND "can_index"
UNION ALL
SELECT "role_id", 'reviews.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'reviews' AND "can_create"
UNION ALL
SELECT "role_id", 'reviews.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'reviews' AND "can_update"
UNION ALL
SELECT "role_id", 'reviews.delete', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'reviews' AND "can_destroy"
UNION ALL
SELECT "role_id", 'leads.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'leads' AND "can_index"
UNION ALL
SELECT "role_id", 'messages.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'messages' AND "can_index";

-- The matrix no longer decides these resources; leaving the rows would show a
-- second, dead copy of the right.
DELETE FROM "staff_member_role_permissions" WHERE "resource" IN ('banners', 'reviews', 'leads', 'messages');
