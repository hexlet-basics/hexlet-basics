class ReviewCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Review

    typelize model: :string
    typelize relations: "Record<string, string>"
    typelize states: "{ key: string, value: string }[]"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
    attribute(:states) { it.class.enum_as_hashes(:states) }
  end

  typelize_from Review

  # one :user, resource: UserResource
  has_one :user
  has_one :language

  attributes :id, :body, :state, :first_name, :last_name, :language_id, :user_id, :pinned

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
