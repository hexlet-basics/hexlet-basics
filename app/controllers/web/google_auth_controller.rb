class Web::GoogleAuthController < Web::ApplicationController
  before_action :validate_google_csrf

  def one_tap
    payload = ApplicationContainer[:google_one_tap].verify_oidc(params[:credential], aud: configus.google.client.id)
    email = payload["email"]
    existing_user = User.find_by(email: email)
    user = GoogleAuthService.authenticate_user(payload)

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
  rescue Google::Auth::IDTokens::SignatureError, Google::Auth::IDTokens::AudienceMismatchError
    f(:error)
    redirect_to root_path
  end

  private

  def validate_google_csrf
    if cookies["g_csrf_token"].blank? || params["g_csrf_token"].blank? ||
       cookies["g_csrf_token"] != params["g_csrf_token"]
      f(:error)
      redirect_to root_path
    end
  end
end
