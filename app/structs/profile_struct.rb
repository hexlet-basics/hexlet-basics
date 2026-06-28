# typed: strict
# frozen_string_literal: true

class ProfileStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :first_name, T.nilable(String)
  const :last_name, T.nilable(String)
  const :nickname, T.nilable(String)
end
