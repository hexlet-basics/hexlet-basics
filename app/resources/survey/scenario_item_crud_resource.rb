class Survey::ScenarioItemCrudResource < ApplicationResource
  typelize_from Survey::Scenario::Item

  attributes :id, :survey_id, :order
  has_one :survey, resource: SurveyResource

  # typelize :state, nullabe: false

  # typelize :string, nullable: true
  # attribute :value_for_select do |obj|
  #   obj.survey.question
  # end
end
