class UserSignUpFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User::SignUpForm

  attributes :first_name, :email, :password
end
