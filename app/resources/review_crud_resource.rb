class ReviewCrudResource < ApplicationResource
  typelize_from Review
  root_key :data

  # one :user, resource: UserResource
  has_one :user
  has_one :language

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :pinned

  typelize_meta(meta: "{ modelName: string, states: { key: string, value: string }[] }")
  meta do
    {
      modelName: object.class.superclass.form_key,
      states: object.class.enum_as_hashes(:states)
    }
  end
end
