# typed: strict
# frozen_string_literal: true

class ReviewStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :user_id, T.nilable(Integer)
  const :language_id, T.nilable(Integer)
  const :pinned, T.nilable(T::Boolean)
  const :state, T.nilable(String)
  const :body, T.nilable(String)
  const :first_name, T.nilable(String)
  const :last_name, T.nilable(String)
end
