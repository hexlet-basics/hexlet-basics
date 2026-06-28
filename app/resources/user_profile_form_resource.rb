# typed: true

class UserProfileFormResource < ApplicationResource
  typelize_from User

  attributes :first_name, :last_name
end
