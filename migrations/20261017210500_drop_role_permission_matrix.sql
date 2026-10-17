-- ADR-0020: every resource of the create/read/update/delete matrix has been
-- carried into grants by the migrations before this one, and nothing reads the
-- table any more. Rows left over name a resource that no operation uses (the
-- earlier migrations delete the ones they convert), so dropping the table loses
-- no access. It is dropped rather than kept empty so the roles screen and the
-- policy cannot grow a second copy of a right.
DROP TABLE "staff_member_role_permissions";
