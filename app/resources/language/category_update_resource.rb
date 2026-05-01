class Language::CategoryUpdateResource < ApplicationResource
  typelize_from Language::Category

  attributes :id, :slug, :name, :header, :description

  typelize slug: :string
  typelize name: :string
  typelize header: :string

  typelize language_landing_page_ids: "number[]"
  attribute :language_landing_page_ids do
    it.items.pluck(:language_landing_page_id)
  end
end
