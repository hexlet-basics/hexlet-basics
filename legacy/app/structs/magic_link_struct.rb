# typed: strict
# frozen_string_literal: true

class MagicLinkStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :email, T.nilable(String)

  validates :email, presence: true
end
