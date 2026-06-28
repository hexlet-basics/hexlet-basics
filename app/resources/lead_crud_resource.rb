# typed: true

class LeadCrudResource < ApplicationResource
  typelize_from LeadStruct

  attributes :contact_method, :contact_value, :ym_client_id

  typelize contact_method: [ enum: [ "telegram", "phone", "whatsapp" ] ],
    contact_value: :string
  typelize ym_client_id: [ :string, nullable: true ]
end
