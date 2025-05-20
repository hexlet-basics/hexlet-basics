class SignInFormResource < ApplicationResource
  typelize_from SignInForm
  root_key :user_sign_in_form

  attributes :email, :password
  typelize email: :string, password: :string

  typelize '"sign_in_form"', nullable: false
  attribute :type do |user|
    "sign_in_form"
  end

  meta do
    {}
  end
end
