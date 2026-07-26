# typed: strict
# frozen_string_literal: true

class SignInStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :email, T.nilable(String)
  const :password, T.nilable(String)

  validates :email, presence: true
  validates :password, presence: true
end
