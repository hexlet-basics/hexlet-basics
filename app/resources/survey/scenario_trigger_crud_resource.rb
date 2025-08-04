class Survey::ScenarioTriggerCrudResource < ApplicationResource
  typelize_from Survey::Scenario::Trigger
  # root_key :data

  attributes :id, :event_name, :event_threshold_count

  typelize id: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do |category|
    false
  end

  # typelize :state, nullabe: false

  # typelize :string, nullable: true
  # attribute :value_for_select do |obj|
  #   obj.survey.question
  # end
  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.form_key
    }
  end
end
