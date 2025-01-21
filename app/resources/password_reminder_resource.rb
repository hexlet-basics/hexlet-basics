class PasswordReminderResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from RemindPasswordForm

  attributes :email
end
