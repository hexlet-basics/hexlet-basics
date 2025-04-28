class Survey::ItemCrudResource < ApplicationResource
  typelize_from Survey::Item

  attributes :id, :survey_id, :value, :state, :order

  typelize :state, nullabe: false

  typelize :string, nullable: true
  attribute :value_for_select do |obj|
    "#{obj.value} (#{obj.survey.slug})"
  end
end
