# typed: true

class SignInFormResource < ApplicationResource
  typelize_from SignInForm

  attributes :email, :password
  typelize email: :string, password: :string
end
