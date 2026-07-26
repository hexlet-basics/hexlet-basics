# typed: strict
# frozen_string_literal: true

class StaffRoleCreateStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :name, T.nilable(String)
  const :description, T.nilable(String)

  validates :name, presence: true
end
