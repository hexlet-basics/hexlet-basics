class SignInFormResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from SignInForm

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { User.to_s.underscore }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from SignInForm

  attributes :email, :password
  typelize email: :string, password: :string

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
