# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    user = User.find_or_initialize_by(email: auth.info.email)
    user.save!

    account = user.accounts.find_or_create_by(provider: auth.provider)
    account.uid = auth.uid
    account.save!

    user
  end
end
