# typed: strict

class Web::GoogleAuthController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :validate_google_csrf

  sig { returns(T.untyped) }
  def one_tap
    payload = T.unsafe(ApplicationContainer)[:google_one_tap].verify_oidc(params[:credential], aud: configus.google.client.id)
    email = payload["email"]
    existing_user = User.find_by(email: email)
    user = GoogleAuthService.authenticate_user(payload)

    if user.persisted?
      sign_in(user)
      f(:success)
      unless existing_user
        signed_up_event_data = {
          user_id: user.id,
          email: T.must(user.email),
          first_name: user.first_name,
          last_name: user.last_name,
          locale: I18n.locale
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

  sig { returns(T.untyped) }
  def validate_google_csrf
    if cookies["g_csrf_token"].blank? || params["g_csrf_token"].blank? ||
       cookies["g_csrf_token"] != params["g_csrf_token"]
      f(:error)
      redirect_to root_path
    end
  end
end
