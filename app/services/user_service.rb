# frozen_string_literal: true

class UserService
  def self.reset_password!(user, suffix)
    token = SecureHelper.generate_token
    user.update!(reset_password_token: token)

    UserMailer.with(user:, suffix:).reset_password.deliver_later
  end
end
