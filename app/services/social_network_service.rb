# typed: true
# frozen_string_literal: true

class SocialNetworkService < ApplicationService
  sig { params(auth: T.untyped).returns(Typed::Success[User]) }
  def self.authenticate_user(auth)
    provider = auth[:provider]
    uid = auth[:uid]
    email = auth[:info][:email].downcase

    # A social identity is keyed by (provider, uid). Looking it up first also
    # handles a changed provider email without creating a duplicate account.
    existing_account = User::Account.find_by(provider: provider, uid: uid)
    user = if existing_account
      # user_id is NOT NULL, so an existing account always has a user
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
          locale: I18n.locale.to_s
        }
        event = UserSignedUpEvent.new(data: signed_up_event_data)
        EventSender.publish_event(event, user)
      end

      user.accounts.find_or_create_by!(provider: provider, uid: uid)
    end

    success_with(user)
  end
end
