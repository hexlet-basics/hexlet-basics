class LeadCrudResource < ApplicationResource
  class MetaResource < ApplicationResource
    typelize_from LeadForm

    typelize model: :string
    typelize relations: "Record<string, string>"

    attribute(:model) { it.class.superclass.form_key }
    attribute(:relations) do
      it.class.respond_to?(:nested_attributes_mapping) ? it.class.nested_attributes_mapping : {}
    end
  end

  typelize_from LeadForm

  attributes :contact_method, :contact_value, :ym_client_id

  typelize contact_method: [ enum: [ "telegram", "phone", "whatsapp" ] ],
    contact_value: :string

  has_one :meta, source: proc { |_params| self }, resource: MetaResource
end
