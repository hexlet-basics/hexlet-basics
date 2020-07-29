# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  def callback
    user = SocialNetworkService.authenticate_user(auth)

    if user.persisted?
      sign_in user
      redirect_to root_path
    else
      redirect_to new_user_path
    end
  end

  private

  def auth
    request.env['omniauth.auth']
  end
end
