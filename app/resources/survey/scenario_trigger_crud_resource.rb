class Survey::ScenarioTriggerCrudResource < ApplicationResource
  typelize_from Survey::Scenario::Trigger

  attributes :id, :event_name, :event_threshold_count

  # typelize :state, nullabe: false

  # typelize :string, nullable: true
  # attribute :value_for_select do |obj|
  #   obj.survey.question
  # end
end
