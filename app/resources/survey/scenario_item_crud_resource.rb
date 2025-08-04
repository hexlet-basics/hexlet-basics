class Survey::ScenarioItemCrudResource < ApplicationResource
  typelize_from Survey::Scenario::Item
  # root_key :data

  attributes :id, :survey_id, :order
  # has_one :survey, resource: SurveyResource

  typelize id: [ :number, nullable: true ]
  typelize survey_id: [ :number, nullable: true ]

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
