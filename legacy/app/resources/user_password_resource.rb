# typed: strict

class UserPasswordResource < ApplicationResource
  typelize_from User

  attributes :password

  typelize password: [ :string, nullable: true ]
end
