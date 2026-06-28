# typed: strict
# frozen_string_literal: true

class GoogleAuthService
  class << self
    extend T::Sig

    sig { params(payload: T.untyped).returns(User::SocialSignupForm) }
    def authenticate_user(payload)
      user = User::SocialSignupForm.find_or_initialize_by(email: payload["email"])
      user.save!

      account = user.accounts.find_or_initialize_by(provider: "google")
      account.uid = payload["sub"]
      account.save!

      user
    end
  end
end
