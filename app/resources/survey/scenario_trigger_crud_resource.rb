class Survey::ScenarioTriggerCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Survey::Scenario::Trigger

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from Survey::Scenario::Trigger
  # root_key :data

  attributes :id, :event_name, :event_threshold_count

  typelize id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do
    false
  end

  # typelize :state, nullable: false

  # typelize :string, nullable: true
  # attribute :value_for_select do |obj|
  #   obj.survey.question
  # end
  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
