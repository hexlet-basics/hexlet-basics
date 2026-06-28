# typed: strict
# frozen_string_literal: true

class PhoneRequestStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :phone, T.nilable(String)

  validates :phone, presence: true
end
