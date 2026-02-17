class PasswordReminderFormResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from RemindPasswordForm

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.to_s.underscore }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from RemindPasswordForm

  attributes :email
  typelize email: :string

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
