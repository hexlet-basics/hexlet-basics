# typed: strict
# frozen_string_literal: true

class StaffMemberService < ApplicationService
  class Payload < T::Struct
    const :staff_member, StaffMember
  end

  class << self
    extend T::Sig

    sig { params(struct: StaffMemberCreateStruct).returns(Typed::Result[Payload, Payload]) }
    def create(struct)
      staff_member = StaffMember.new(**struct.attributes)
      payload = Payload.new(staff_member:)
      return fail_with(payload) unless staff_member.save

      success_with(payload)
    end

    sig { params(id: T.untyped, struct: StaffMemberUpdateStruct).returns(Typed::Result[Payload, Payload]) }
    def update(id, struct)
      staff_member = StaffMember.find(id)
      payload = Payload.new(staff_member:)
      return fail_with(payload) unless staff_member.update(struct.attributes)

      success_with(payload)
    end
  end
end
