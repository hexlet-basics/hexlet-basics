class Language::CategoryItemCrudResource < ApplicationResource
  typelize_from Language::Category::Item
  # root_key :data

  attributes :id, :language_category_id, :language_landing_page_id

  typelize id: [ :number, nullable: true ]
  typelize language_category_id: [ :number, nullable: true ]
  typelize language_landing_page_id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do |category|
    false
  end
end
