class UserSignUpFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User::SignUpForm
  root_key :user_sign_up_form

  attributes :first_name, :email, :password

  meta do
    {}
  end
end
