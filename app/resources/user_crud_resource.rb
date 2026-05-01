class UserCrudResource < ApplicationResource
  typelize_from Admin::UserForm

  attributes :id, :email, :first_name, :last_name, :admin
end
