class Survey::ScenarioCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Survey::Scenario

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from Survey::Scenario

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

  # typelize :state, nullable: false

  typelize :string, nullable: true
  attribute :survey_item_value do
    it.survey_item&.value
  end

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
