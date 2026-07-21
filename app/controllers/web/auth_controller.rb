# typed: strict
# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  allow_unauthenticated_access

  sig { void }
  def callback
    result = SocialNetworkService.authenticate_user(auth)
    user = result.payload

    sign_in(user)
    f(:success)
    # js_event_options = {
    #   user: result.user
    # }
    # js_event(result.is_new ? :signed_up : :signed_in, js_event_options)
    redirect_to root_path
    # TODO Посылать расширенный контекст в sentry
    # rescue ActiveRecord::RecordInvalid => e
    # e.rollbar_context = { auth_hash: auth }
    # raise
  end

  private

  sig { void }
  def auth
    request.env["omniauth.auth"]
  end
end
