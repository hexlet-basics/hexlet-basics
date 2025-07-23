class Language::CategoryQnaItemCrudResource < ApplicationResource
  typelize_from Language::Category::QnaItem
  # root_key :data

  attributes :id, :question, :answer

  typelize id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do |category|
    false
  end

  # typelize_meta meta: "{ modelName: string }"
  # meta do
  #   {
  #     modelName: object.class.superclass.to_s.underscore
  #   }
  # end
end
