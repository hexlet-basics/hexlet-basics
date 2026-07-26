# typed: strict

class StaffRoleCrudResource < ApplicationResource
  typelize_from StaffMember::Role

  attributes :id, :name, :description
  has_many :permissions, resource: StaffRolePermissionResource
end
