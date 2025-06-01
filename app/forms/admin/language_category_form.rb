class Admin::LanguageCategoryForm < Language::Category
  include ActiveFormModel

  permit :name, :slug, :header, :description,
    items_attributes: [ :id, :language_landing_page_id, :_destroy ]

  accepts_nested_attributes_for :items, allow_destroy: true
end
