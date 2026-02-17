class UserCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from Admin::UserForm

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from Admin::UserForm

  attributes :id, :email, :first_name, :last_name, :admin

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
