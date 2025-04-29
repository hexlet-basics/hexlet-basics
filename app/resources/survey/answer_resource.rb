class Survey::AnswerResource < ApplicationResource
  typelize_from Survey::Answer

  attributes :id, :survey_id, :survey_item_id, :user_id

  typelize :state, nullabe: false

  typelize :string, nullable: true
  attribute :survey_slug do |obj|
    obj.survey.slug
  end

  typelize :string, nullable: true
  attribute :survey_item_value do |obj|
    obj.survey_item&.value
  end
end
