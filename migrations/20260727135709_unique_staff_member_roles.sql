-- Enforce the legacy StaffMember::Role name uniqueness as a real DB constraint
-- (parity with legacy/app/models/staff_member/role.rb `validates :name,
-- uniqueness: true`). The handler does not pre-check: a duplicate surfaces as
-- ent.IsConstraintError and the central ogen ErrorHandler maps it to 409.
--
-- PRE-DEPLOY: added to a table holding real legacy data; verify no duplicates
-- exist on prod before shipping, or the CREATE UNIQUE INDEX will fail:
--   SELECT name, count(*) FROM staff_member_roles GROUP BY name HAVING count(*) > 1;
CREATE UNIQUE INDEX "staff_member_roles_name_uniq" ON "staff_member_roles" ("name");
