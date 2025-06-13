class ReviewCrudResource < ApplicationResource
  typelize_from Review
  root_key :review

  # one :user, resource: UserResource
  has_one :user
  has_one :language

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :pinned

  typelize_meta(meta: "{ states: { key: string, value: string }[] }")
  meta do
    { states: object.class.enum_as_hashes(:states) }
  end
end
