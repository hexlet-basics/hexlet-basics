# typed: strict

class StaffRolePermissionResource < ApplicationResource
  typelize_from StaffMember::Role::Permission

  attributes :id, :role_id, :resource, :can_index, :can_create, :can_update, :can_destroy
end
