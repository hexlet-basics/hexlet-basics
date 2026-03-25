# frozen_string_literal: true

class UserService
  def self.reset_password!(user, suffix)
    UserMailer.with(user:, suffix:).reset_password.deliver_later
  end
end
