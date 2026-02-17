class Survey::ScenarioItemCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Survey::Scenario::Item

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from Survey::Scenario::Item
  # root_key :data

  attributes :id, :survey_id, :order
  # has_one :survey, resource: SurveyResource

  typelize id: [ :number, nullable: true ]
  typelize survey_id: [ :number, nullable: true ]

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
