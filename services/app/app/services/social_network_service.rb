# frozen_string_literal: true

class SocialNetworkService
  def self.authenticate_user(auth)
    user = User::SocialSignupForm.find_or_initialize_by(email: auth.info.email)
    user.save!

    account = user.accounts.find_or_initialize_by(provider: auth.provider)
    account.uid = auth.uid
    account.save!

    user
  # NOTE Added to catch an invalid record error during auth via github
  rescue ActiveRecord::RecordInvalid => e
    e.rollbar_context = {
      auth_response: auth,
      info: auth.info
    }
    raise
  end
end
