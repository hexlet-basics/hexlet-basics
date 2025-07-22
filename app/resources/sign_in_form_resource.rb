class SignInFormResource < ApplicationResource
  typelize_from SignInForm
  root_key :data

  attributes :email, :password
  typelize email: :string, password: :string

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: User.to_s.underscore
    }
  end
end
