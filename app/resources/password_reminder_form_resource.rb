class PasswordReminderFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from RemindPasswordForm

  attributes :email

  meta do
    {}
  end
end
