# typed: strict
# frozen_string_literal: true

class StaffMemberUpdateStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :role_id, T.nilable(Integer)
  const :allowed_locales, T.nilable(T::Array[String])

  validates :role_id, presence: true
end
