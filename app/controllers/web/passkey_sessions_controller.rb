# typed: strict
# frozen_string_literal: true

class Web::PasskeySessionsController < Web::ApplicationController
  allow_unauthenticated_access
  before_action :redirect_if_authenticated

  rescue_from WebAuthn::Error, ActiveRecord::RecordNotFound, with: :handle_failure

  # Returns WebAuthn request options (JSON) and stashes the challenge in the session.
  # Discoverable credentials → no allow list, so the user picks an account.
  sig { returns(T.untyped) }
  def new
    options = WebAuthn::Credential.options_for_get(user_verification: "preferred")
    session[:passkey_authentication_challenge] = options.challenge

    render json: options
  end

  sig { returns(T.untyped) }
  def create
    webauthn_credential = WebAuthn::Credential.from_get(JSON.parse(params.require(:credential)))
    stored = User::Credential.find_by!(external_id: webauthn_credential.id)

    webauthn_credential.verify(
      session.delete(:passkey_authentication_challenge),
      public_key: stored.public_key,
      sign_count: stored.sign_count
    )
    stored.update!(sign_count: webauthn_credential.sign_count)

    user = stored.user
    publish_signed_in_event(user)
    sign_in(user)

    f(:success)
    redirect_to after_authentication_url
  end

  private

  sig { params(user: User).void }
  def publish_signed_in_event(user)
    event = UserSignedInEvent.new(data: {
      user_id: user.id, occurrence_count: -1, email: user.email, locale: I18n.locale.to_s
    })
    publish_event(event, user)
    js_event(event)
  end

  sig { returns(T.untyped) }
  def handle_failure
    f(:error, type: :alert)
    redirect_to new_session_path
  end
end
