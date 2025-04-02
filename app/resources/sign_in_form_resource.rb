class SignInFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from SignInForm

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
