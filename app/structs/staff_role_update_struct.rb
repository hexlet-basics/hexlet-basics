# typed: strict
# frozen_string_literal: true

class StaffRoleUpdateStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :name, T.nilable(String)
  const :description, T.nilable(String)

  validates :name, presence: true
end
