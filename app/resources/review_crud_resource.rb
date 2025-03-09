class ReviewCrudResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from Review
  root_key :review

  # one :user, resource: UserResource
  has_one :user
  has_one :language

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id

  typelize_meta meta: "{ state_events: Array<[string, string]>}"
  meta do
    { state_events: object.aasm.events_for_select }
  end
end
