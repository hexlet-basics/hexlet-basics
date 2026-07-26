# typed: strict
# frozen_string_literal: true

class RemindPasswordStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :email, T.nilable(String)

  validates :email, presence: true
end
