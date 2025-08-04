class UserCrudResource < ApplicationResource
  typelize_from Admin::UserForm
  root_key :data

  attributes :id, :email, :first_name, :last_name, :admin

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.superclass.form_key
    }
  end
end
