# frozen_string_literal: true

class Web::AuthController < Web::ApplicationController
  def callback
    result = SocialNetworkService.authenticate_user(auth)

    if result.fail?
      f(:error)
      return redirect_to(new_session_path)
    end

    if result.is_new
      sign_up(result.user)
      sign_in(result.user)
      fill_guests_data(result.user)
    else
      event_store.within { sign_in result.user }
                 .subscribe(to: UserSignedInEvent) { |event| event_to_js(event) }
                 .call
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
