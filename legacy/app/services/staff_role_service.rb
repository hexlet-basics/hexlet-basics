# typed: strict
# frozen_string_literal: true

class StaffRoleService < ApplicationService
  class Payload < T::Struct
    const :role, StaffMember::Role
  end

  class << self
    extend T::Sig

    sig { params(struct: StaffRoleCreateStruct).returns(Typed::Result[Payload, Payload]) }
    def create(struct)
      role = StaffRoleMutator.create(struct)
      payload = Payload.new(role:)
      return fail_with(payload) unless role.persisted?

      success_with(payload)
    end

    sig { params(id: String, struct: StaffRoleUpdateStruct).returns(Typed::Result[Payload, Payload]) }
    def update(id, struct)
      role = StaffMember::Role.find(id)
      payload = Payload.new(role:)
      return fail_with(payload) unless StaffRoleMutator.update(role, struct)

      success_with(payload)
    end
  end
end
