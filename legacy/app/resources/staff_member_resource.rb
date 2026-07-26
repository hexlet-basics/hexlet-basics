# typed: strict

class StaffMemberResource < ApplicationResource
  typelize_from StaffMember

  has_one :user, resource: UserResource
  has_one :role, resource: StaffRoleResource

  typelize id: :number
  typelize user_id: :number
  typelize role_id: :number
  attributes :id, :user_id, :role_id, :allowed_locales

  typelize allowed_locales: "string[]"
end
