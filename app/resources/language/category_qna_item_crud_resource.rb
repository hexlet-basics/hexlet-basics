class Language::CategoryQnaItemCrudResource < ApplicationResource
  typelize_from Language::Category::QnaItem
  # root_key :data

  attributes :id, :question, :answer

  typelize id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do
    false
  end
end
