class SurveyCrudResource < ApplicationResource
  typelize_from Survey

  class MetaResource < ApplicationResource
    typelize_from Survey

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize item_states: "Record<string, unknown>[]"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:item_states) { Survey::Item.states }
  end

  # one :user, resource: UserResource
  # has_one :user
  # has_one :language
  has_many :items, resource: Survey::ItemCrudResource, key: "items_attributes"

  attributes :id,
    :state,
    :question,
    :description,
    :slug

  # typelize :state, nullable: false

  # typelize :string, nullable: true
  # attribute :parent_survey_item_value do |obj|
  #   obj.parent_survey_item&.value
  # end

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
