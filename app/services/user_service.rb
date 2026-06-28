# typed: strict
# frozen_string_literal: true

class UserService
  extend T::Sig

  sig { params(user: User, suffix: T.untyped).void }
  def self.reset_password!(user, suffix)
    UserMailer.with(user:, suffix:).reset_password.deliver_later
  end
end
