# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  include EventEmitterConcern

  def callback
    result = SocialNetworkService.authenticate_user(auth)

    if result.fail?
      f(:error)
      return redirect_to(new_session_path)
    end

    if result.is_new
      sign_up_event(result.user)
      sign_in(result.user)
      sign_up_progress_events(result.user)
    else
      sign_in(result.user)
    end
    f(:success)
    redirect_to root_path
  end

  def failure
    f(:error, values: { text: params[:message] })
    redirect_to(new_session_path)
  end

  private

  def auth
    request.env["omniauth.auth"]
  end
end
