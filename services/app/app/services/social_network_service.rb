# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    user = User::SocialSignupForm.find_or_initialize_by(email: auth.info.email)
    user.locale = I18n.locale
    user.valid?
    user.save!

    account = user.accounts.find_or_initialize_by(provider: auth.provider)
    account.uid = auth.uid
    account.save!

    user
  end
end
