class SurveyCrudResource < ApplicationResource
  typelize_from Survey
  root_key :survey

  # one :user, resource: UserResource
  # has_one :user
  # has_one :language
  has_many :items, resource: Survey::ItemResource

  attributes :id, :state, :question, :description, :slug

  typelize :state, nullabe: false

  typelize_meta meta: "{ item_states: Record<string, unknown>[]}"
  meta do
    { item_states: Survey::Item.enum_as_hashes(:states) }
  end
end
