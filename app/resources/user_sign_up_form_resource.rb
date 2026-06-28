# typed: true

class UserSignUpFormResource < ApplicationResource
  typelize_from User

  attributes :first_name, :email, :password

  typelize password: [ :string, nullable: true ]
end
