class UserProfileFormResource < ApplicationResource
  typelize_from User::ProfileForm
  root_key :data

  attributes :first_name, :last_name

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.to_s.underscore
    }
  end
end
