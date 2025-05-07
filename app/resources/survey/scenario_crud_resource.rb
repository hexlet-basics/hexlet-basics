class Survey::ScenarioCrudResource < ApplicationResource
  typelize_from Survey::Scenario
  root_key :survey_scenario

  has_many :items, resource: Survey::ScenarioItemCrudResource
  has_many :triggers, resource: Survey::ScenarioTriggerCrudResource
  has_one :survey_item, resource: Survey::ItemCrudResource

  attributes :id,
    :state,
    :name,
    # :event_name,
    # :event_threshold_count,
    :survey_item_id

  # typelize :state, nullabe: false

  typelize :string, nullable: true
  attribute :survey_item_value do |obj|
    obj.survey_item&.value
  end

  # typelize_meta meta: "{ item_states: Record<string, unknown>[]}"
  meta do
    {}
    # { item_states: Survey::Item.enum_as_hashes(:states) }
  end
end
