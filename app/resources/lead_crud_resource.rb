class LeadCrudResource < ApplicationResource
  typelize_from Lead
  root_key :lead

  attributes :contact_method, :contact_value
end
