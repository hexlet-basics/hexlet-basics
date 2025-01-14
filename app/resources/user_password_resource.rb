class UserPasswordResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from User::PasswordForm

  attributes :password
end
