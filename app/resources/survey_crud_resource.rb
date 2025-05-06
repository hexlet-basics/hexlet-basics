class SurveyCrudResource < ApplicationResource
  typelize_from Survey
  root_key :survey

  # one :user, resource: UserResource
  # has_one :user
  # has_one :language
  has_many :items, resource: Survey::ItemResource

  attributes :id,
    :state,
    :question,
    :description,
    :slug,
    :parent_survey_item_id,
    :parent_survey_id,
    :run_always,
    :run_after_finishing_lessons_count

  typelize :state, nullabe: false

  typelize :string, nullable: true
  attribute :parent_survey_item_value do |obj|
    obj.parent_survey_item&.value
  end

  typelize_meta meta: "{ item_states: Record<string, unknown>[]}"
  meta do
    { item_states: Survey::Item.enum_as_hashes(:states) }
  end
end
