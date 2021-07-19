# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  def callback
    email = auth[:info][:email].downcase
    existing_user = User.find_by(email: email)

    user = SocialNetworkService.authenticate_user(auth)

    if user.persisted?
      sign_in user
      f(:success)
      js_event_options = {
        user: user
      }
      js_event(:signed_up, js_event_options) unless existing_user
      redirect_to root_path
    else
      redirect_to new_user_path
    end
  rescue ActiveRecord::RecordInvalid => e
    e.rollbar_context = { auth_hash: auth }
    raise
  end

  private

  def auth
    request.env['omniauth.auth']
  end
end
