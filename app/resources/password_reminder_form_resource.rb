class PasswordReminderFormResource
  include Alba::Resource
  include Typelizer::DSL

  typelize_from RemindPasswordForm

  attributes :email

  typelize email: :string

  meta do
    {}
  end
end
