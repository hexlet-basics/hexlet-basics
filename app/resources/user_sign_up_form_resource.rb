class UserSignUpFormResource < ApplicationResource
  typelize_from User::SignUpForm
  root_key :data

  attributes :first_name, :email, :password

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.to_s.underscore
    }
  end
end
