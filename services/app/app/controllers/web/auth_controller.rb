# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  def callback
    user = SocialNetworkService.authenticate_user(auth)

    if user.persisted?
      sign_in user
      f(:success)
      redirect_to root_path
    else
      redirect_to new_user_path
    end
  end

  private

  def auth
    # raise request[env.inspect
    request.env['omniauth.auth']
  end
end
