class PasswordReminderFormResource < ApplicationResource
  typelize_from RemindPasswordForm
  root_key :data

  attributes :email
  typelize email: :string

  typelize_meta meta: "{ modelName: string }"
  meta do
    {
      modelName: object.class.to_s.underscore
    }
  end
end
