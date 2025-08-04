class Survey::ItemCrudResource < ApplicationResource
  typelize_from Survey::Item
  # root_key :data

  attributes :id, :survey_id, :value, :state, :order

  typelize id: [ :number, nullable: true ]
  typelize survey_id: [ :number, nullable: true ]
  typelize order: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do |category|
    false
  end

  typelize :string, nullable: true
  attribute :value_for_select do |obj|
    "#{obj.value} (#{obj.survey.slug})"
  end

  typelize :string, nullable: true
  attribute :tag_list do |obj|
    obj.tag_list.join ", "
  end

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.form_key
    }
  end
end
