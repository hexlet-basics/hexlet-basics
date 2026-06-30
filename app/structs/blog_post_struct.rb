# typed: strict
# frozen_string_literal: true

class BlogPostStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :creator_id, T.nilable(Integer)
  const :language_id, T.nilable(Integer)
  const :rich_body, T.nilable(String)
  const :slug, T.nilable(String)
  const :name, T.nilable(String)
  const :description, T.nilable(String)
  const :state, T.nilable(String)
end
