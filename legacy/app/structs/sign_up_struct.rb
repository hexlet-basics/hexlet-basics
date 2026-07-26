# typed: strict
# frozen_string_literal: true

class SignUpStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :email, T.nilable(String)
  const :password, T.nilable(String)
  const :first_name, T.nilable(String)

  validates :password, presence: true, length: { minimum: 6 }
end
