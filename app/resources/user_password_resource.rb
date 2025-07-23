class UserPasswordResource < ApplicationResource
  typelize_from User::PasswordForm
  root_key :data

  attributes :password

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.to_s.underscore
    }
  end
end
