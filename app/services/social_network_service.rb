# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    existing_account = User::Account.find_by(uid: auth[:uid])
    email = auth[:info][:email].downcase
    auth_user = if existing_account
                  existing_account.user
                else
                  User::SocialSignupForm.find_or_initialize_by(email: email)
                end

    is_new = false

    if auth_user.new_record?
      auth_user.save!
      is_new = true
    else
      auth_user.email = email
      auth_user.save!
    end

    account = auth_user.accounts.find_or_initialize_by(provider: auth[:provider])
    account.uid = auth[:uid]
    account.save!

    { user: auth_user, is_new: is_new }
  end
end
