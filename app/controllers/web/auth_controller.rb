# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  def callback
    result = SocialNetworkService.authenticate_user(auth)

    sign_in result.user
    f(:success)
    js_event_options = {
      user: result.user
    }
    js_event(result.is_new ? :signed_up : :signed_in, js_event_options)
    redirect_to root_path
    # TODO Посылать расширенный контекст в sentry
    # rescue ActiveRecord::RecordInvalid => e
    # e.rollbar_context = { auth_hash: auth }
    # raise
  end

  private

  def auth
    request.env["omniauth.auth"]
  end
end
