# typed: strict
# frozen_string_literal: true

class LanguageCategoryStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :name, T.nilable(String)
  const :slug, T.nilable(String)
  const :header, T.nilable(String)
  const :description, T.nilable(String)
  const :language_landing_page_ids, T.nilable(T::Array[Integer])
end
