# frozen_string_literal: true

class Admin::LanguageCategoryForm < Language::Category
  include ActiveFormModel

  permit :name, :slug, :header, :description
end
