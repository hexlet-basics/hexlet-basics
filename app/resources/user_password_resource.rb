class UserPasswordResource < ApplicationResource
  typelize_from User::PasswordForm

  attributes :password

  typelize password: [ :string, nullable: true ]
end
