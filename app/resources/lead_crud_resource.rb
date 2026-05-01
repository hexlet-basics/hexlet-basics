class LeadCrudResource < ApplicationResource
  typelize_from LeadForm

  attributes :contact_method, :contact_value, :ym_client_id

  typelize contact_method: [ enum: [ "telegram", "phone", "whatsapp" ] ],
    contact_value: :string
end
