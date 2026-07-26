# typed: strict
# frozen_string_literal: true

class GoogleAuthService < ApplicationService
  class << self
    extend T::Sig

    sig { params(payload: T::Hash[String, T.untyped]).returns(Typed::Success[User]) }
    def authenticate_user(payload)
      uid = payload["sub"]

      # A social identity is keyed by (provider, uid). Looking it up first also
      # handles a changed Google email without creating a duplicate account.
      existing_account = User::Account.find_by(provider: "google", uid: uid)
      user = existing_account ? existing_account.user : User.find_or_initialize_by(email: payload["email"])
      user.save! if user.new_record?

      user.accounts.find_or_create_by!(provider: "google", uid: uid)

      success_with(user)
    end
  end
end
