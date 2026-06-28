# typed: true

class UserProfileFormResource < ApplicationResource
  typelize_from User::ProfileForm

  attributes :first_name, :last_name
end
