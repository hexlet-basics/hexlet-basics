# typed: strict

class StaffRoleResource < ApplicationResource
  typelize_from StaffMember::Role

  attributes :id, :name, :description, :created_at

  typelize :number
  attribute :permissions_count do
    it.permissions.count { |p| p.can_index? || p.can_create? || p.can_update? || p.can_destroy? }
  end
end
