class LeadCrudResource < ApplicationResource
  typelize_from Lead
  root_key :data

  attributes :contact_method, :contact_value, :ym_client_id

  # TODO: use enum
  # typelize contact_method: '"telegram" | "telephone" | "whatsapp"',
  #   contact_value: :string

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.form_key
    }
  end
end
