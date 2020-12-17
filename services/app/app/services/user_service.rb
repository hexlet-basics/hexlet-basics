# frozen_string_literal: true

class UserService
  def self.reset_password!(user)
    token = SecureHelper.generate_token
    user.update!(reset_password_token: token)

    UserMailer.reset_password(user.id).deliver_later
  end
end
