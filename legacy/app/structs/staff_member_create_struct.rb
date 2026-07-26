# typed: strict
# frozen_string_literal: true

class StaffMemberCreateStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :user_id, T.nilable(Integer)
  const :role_id, T.nilable(Integer)
  const :allowed_locales, T.nilable(T::Array[String])

  validates :user_id, presence: true
  validates :role_id, presence: true
end
