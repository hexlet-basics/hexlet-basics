class Language::CategoryItemResource < ApplicationResource
  typelize_from Language::Category::Item
  # root_key :user

  attributes :id, :language_category_id, :language_landing_page_id, :_destroy
end
