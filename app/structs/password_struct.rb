# typed: strict
# frozen_string_literal: true

class PasswordStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :password, T.nilable(String)

  validates :password, presence: true, length: { minimum: 6 }
end
