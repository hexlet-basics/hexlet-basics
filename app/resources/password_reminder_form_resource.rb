class PasswordReminderFormResource < ApplicationResource
  typelize_from RemindPasswordForm

  attributes :email
  typelize email: :string
end
