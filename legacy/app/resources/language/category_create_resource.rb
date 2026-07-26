# typed: strict

class Language::CategoryCreateResource < ApplicationResource
  typelize_from Language::Category

  typelize id: [ :number, nullable: true ]
  attributes :id,
    slug: [ String, true ],
    name: [ String, true ],
    header: [ String, true ],
    description: [ String, true ]

  typelize language_landing_page_ids: "number[]"
  attribute :language_landing_page_ids do
    it.items.pluck(:language_landing_page_id)
  end
end
