# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    existing_account = User::Account.find_by(uid: auth[:uid])
    email = auth.dig(:info, :email)&.downcase
    user = if existing_account
      existing_account.user
    else
      User.find_or_initialize_by(email:)
    end
    is_new = user.new_record?
    ActiveRecord::Base.transaction do
      user.save! if is_new

      account = user.accounts.find_or_initialize_by(provider: auth[:provider])
      account.uid = auth[:uid]
      account.save!
    end
    ServiceResult.success user: user, is_new:
  rescue ActiveRecord::RecordInvalid => _e
    ServiceResult.error
  end
end
