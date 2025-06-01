class Language::CategoryCrudResource < ApplicationResource
  typelize_from Language::Category
  root_key :language_category

  attributes :id, :slug, :name, :header, :description

  has_many :items, resource: Language::CategoryItemResource
  has_many :qna_items, resource: Language::CategoryQnaItemCrudResource

  typelize_meta meta: "{ landingPagesForCategories: LanguageCategoryCrudData[] }"
  meta do
    landing_pages_for_categories = Language::LandingPage.web
      .where(listed: true)
      .merge(Language.ordered)

    {
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages_for_categories)
    }
  end
end
