# typed: strict
# frozen_string_literal: true

class UserStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :email, T.nilable(String)
  const :first_name, T.nilable(String)
  const :last_name, T.nilable(String)
  const :admin, T.nilable(T::Boolean)
end
