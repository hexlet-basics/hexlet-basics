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
      unless existing_user
        signed_up_event_data = {
          user_id: user.id,
          email: user.email,
          first_name: user.first_name,
          last_name: user.last_name
        }
        js_event(UserSignedUpEvent.new(data: signed_up_event_data))
      end
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
