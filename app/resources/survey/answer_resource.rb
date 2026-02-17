class Survey::AnswerResource < ApplicationResource
  typelize_from Survey::Answer

  attributes :id, :survey_id, :survey_item_id, :user_id, :created_at

  typelize :state, nullable: false

  typelize :string, nullable: true
  attribute :survey_slug do
    it.survey.slug
  end

  typelize :string, nullable: true
  attribute :survey_item_value do
    it.survey_item&.value
  end
end
