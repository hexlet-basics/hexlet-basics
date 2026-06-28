# typed: strict
# frozen_string_literal: true

class LandingPageStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :meta_title, T.nilable(String)
  const :meta_description, T.nilable(String)
  const :listed, T.nilable(T::Boolean)
  const :slug, T.nilable(String)
  const :footer, T.nilable(T::Boolean)
  const :footer_name, T.nilable(String)
  const :order, T.nilable(String)
  const :main, T.nilable(T::Boolean)
  const :state, T.nilable(String)
  const :header, T.nilable(String)
  const :name, T.nilable(String)
  const :description, T.nilable(String)
  const :language_id, T.nilable(Integer)
  const :language_category_id, T.nilable(Integer)
  const :landing_page_to_redirect_id, T.nilable(Integer)
  const :used_in_header, T.nilable(String)
  const :used_in_description, T.nilable(String)
  const :outcomes_header, T.nilable(String)
  const :outcomes_description, T.nilable(String)
end
