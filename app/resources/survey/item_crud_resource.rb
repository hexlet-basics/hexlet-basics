class Survey::ItemCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Survey::Item

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from Survey::Item
  # root_key :data

  attributes :id, :survey_id, :value, :state, :order

  typelize id: [ :number, nullable: true ]
  typelize survey_id: [ :number, nullable: true ]
  typelize order: [ :number, nullable: true ]

  typelize :boolean
  attribute :_destroy do
    false
  end

  typelize :string, nullable: true
  attribute :value_for_select do
    "#{it.value} (#{it.survey.slug})"
  end

  typelize :string, nullable: true
  attribute :tag_list do
    it.tag_list.join ", "
  end

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
