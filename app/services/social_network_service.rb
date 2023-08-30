# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    existing_account = User::Account.find_by(uid: auth[:uid])
    email = auth[:info][:email].downcase
    user = if existing_account
             existing_account.user
           else
             User.find_or_initialize_by(email: email)
           end

    is_new = false

    if user.new_record?
      is_new = true
    else
      user.email = email
    end

    account = user.accounts.find_or_initialize_by(provider: auth[:provider])
    account.uid = auth[:uid]
    ActiveRecord::Base.transaction do
      user.save!
      account.save!
    end
    ServiceResult.success user: user, is_new: is_new
  end
end
