class Survey::ScenarioCrudResource < ApplicationResource
  typelize_from Survey::Scenario
  root_key :data

  has_many :items, proc { |items, params, user|
    items.order(order: :asc)
  }, resource: Survey::ScenarioItemCrudResource, key: "items_attributes"
  has_many :triggers, resource: Survey::ScenarioTriggerCrudResource, key: "triggers_attributes"
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

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.to_s.underscore
    }
  end
end
