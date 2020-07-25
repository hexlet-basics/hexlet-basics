# frozen_string_literal: true

class UserManager
  def self.authenticate_user_by_social_network(auth)
    existing_user = User.find_by(email: auth.info.email)

    if existing_user
      account = existing_user.accounts.find_or_create_by(provider: auth.provider)
      account.uid = auth.uid
      account.save!

      existing_user
    else
      user = User.new(email: auth.info.email)
      user.save!

      account = user.accounts.build(provider: auth.provider, uid: auth.uid)
      account.save!

      user
    end
  end
end
