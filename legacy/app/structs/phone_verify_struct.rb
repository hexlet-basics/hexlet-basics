# typed: strict
# frozen_string_literal: true

class PhoneVerifyStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :phone, T.nilable(String)
  const :code, T.nilable(String)

  validates :phone, presence: true
  validates :code, presence: true
end
