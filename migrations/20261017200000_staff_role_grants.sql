-- ADR-0018: a Staff Role holds Permissions as grants, one row each, instead of
-- create/read/update/delete bits per resource. Blog posts move first; the other
-- resources keep the matrix until their own migration.
CREATE TABLE "staff_member_role_grants" (
  "id" bigserial NOT NULL,
  "created_at" timestamp NOT NULL,
  "permission" character varying NOT NULL,
  "role_id" bigint NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "staff_member_role_grants_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "staff_member_roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE UNIQUE INDEX "index_staff_member_role_grants_on_role_id_and_permission" ON "staff_member_role_grants" ("role_id", "permission");

-- Carry every role's blog post rights over bit for bit, so nobody gains or loses
-- access. The old policy let can_index open the list and the single post,
-- can_create the create, can_update the update, the related-courses set and
-- the AI suggestion, and can_destroy the delete; each bit becomes the grant
-- that opens exactly those operations.
INSERT INTO "staff_member_role_grants" ("role_id", "permission", "created_at", "updated_at")
SELECT "role_id", 'blog_posts.read', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'blog_posts' AND "can_index"
UNION ALL
SELECT "role_id", 'blog_posts.create', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'blog_posts' AND "can_create"
UNION ALL
SELECT "role_id", 'blog_posts.update', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'blog_posts' AND "can_update"
UNION ALL
SELECT "role_id", 'blog_posts.delete', now(), now() FROM "staff_member_role_permissions"
  WHERE "resource" = 'blog_posts' AND "can_destroy";

-- The matrix no longer decides anything for blog posts; leaving the rows would
-- show a second, dead copy of the right on the roles screen.
DELETE FROM "staff_member_role_permissions" WHERE "resource" = 'blog_posts';
