class UserSignUpFormResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from User::SignUpForm

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.to_s.underscore }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from User::SignUpForm

  attributes :first_name, :email, :password

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
