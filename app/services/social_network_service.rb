# frozen_string_literal: true
# typed: true

class SocialNetworkService
  def self.authenticate_user(auth)
    existing_account = User::Account.find_by(uid: auth[:uid])
    email = auth[:info][:email].downcase
    user = if existing_account
      existing_account.user
    else
      User.find_or_initialize_by(email: email)
    end

    ActiveRecord::Base.transaction do
      if user.new_record?
        user.save!
        signed_up_event_data = {
          user_id: user.id,
          email: T.must(user.email),
          first_name: user.first_name,
          last_name: user.last_name,
          locale: I18n.locale
        }
        event = UserSignedUpEvent.new(data: signed_up_event_data)
        publish_event(event, user)
      end

      account = user.accounts.find_or_initialize_by(provider: auth[:provider])
      account.uid = auth[:uid]
      account.save!
    end
    # ServiceResult.success user: user, is_new: is_new
  end
end
