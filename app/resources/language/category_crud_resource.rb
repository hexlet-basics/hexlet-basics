class Language::CategoryCrudResource < ApplicationResource
  typelize_from Language::Category
  root_key :data

  attributes :id, :slug, :name, :header, :description

  many :items, resource: Language::CategoryItemCrudResource # , key: "items_attributes"

  has_many :qna_items, resource: Language::CategoryQnaItemCrudResource # , key: "qna_items_attributes"

  typelize_meta meta: "{ modelName: string, landingPagesForCategories: LanguageCategoryCrudData[] }"
  meta do
    landing_pages_for_categories = Language::LandingPage.web
      # .where(listed: true)
      .merge(Language.ordered)

    {
      modelName: object.class.superclass.to_s.underscore.tr("/", "_"),
      landingPagesForCategories: Language::LandingPageForListsResource.new(landing_pages_for_categories)
    }
  end
end
