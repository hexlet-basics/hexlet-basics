# typed: true

class UserCrudResource < ApplicationResource
  typelize_from User

  attributes :id, :email, :first_name, :last_name, :admin
end
