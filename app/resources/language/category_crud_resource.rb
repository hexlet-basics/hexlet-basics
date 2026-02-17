class Language::CategoryCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Language::Category

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize landing_pages_for_categories: "LanguageLandingPageForLists[]"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:landing_pages_for_categories) do
      landing_pages_for_categories = Language::LandingPage.web
        .merge(Language.ordered)
      Language::LandingPageForListsResource.new(landing_pages_for_categories)
    end
  end

  typelize_from Language::Category

  attributes :id, :slug, :name, :header, :description

  many :items, resource: Language::CategoryItemCrudResource, key: "items_attributes"

  has_many :qna_items, resource: Language::CategoryQnaItemCrudResource, key: "qna_items_attributes"

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
