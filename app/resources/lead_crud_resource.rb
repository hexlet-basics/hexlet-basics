class LeadCrudResource < ApplicationResource
  typelize_from Lead
  root_key :lead

  attributes :contact_method, :contact_value, :ym_client_id

  typelize contact_method: '"telegram" | "telephone" | "whatsapp"',
    contact_value: :string
end
