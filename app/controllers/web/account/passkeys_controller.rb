# typed: strict
# frozen_string_literal: true

class Web::Account::PasskeysController < Web::Account::ApplicationController
  rescue_from WebAuthn::Error, with: :handle_failure

  # Returns WebAuthn creation options (JSON) for the signed-in user and stashes
  # the challenge in the session. Generates a stable webauthn_id on first use.
  sig { void }
  def new
    user = T.must(current_user)
    user.update!(webauthn_id: WebAuthn.generate_user_id) if user.webauthn_id.blank?

    options = WebAuthn::Credential.options_for_create(
      user: { id: user.webauthn_id, name: user.to_s },
      exclude: user.credentials.pluck(:external_id),
      authenticator_selection: { resident_key: "preferred", user_verification: "preferred" }
    )
    session[:passkey_registration_challenge] = options.challenge

    render json: options
  end

  sig { void }
  def create
    user = T.must(current_user)
    webauthn_credential = WebAuthn::Credential.from_create(JSON.parse(params.require(:credential)))
    webauthn_credential.verify(session.delete(:passkey_registration_challenge))

    user.credentials.create!(
      external_id: webauthn_credential.id,
      public_key: webauthn_credential.public_key,
      sign_count: webauthn_credential.sign_count,
      nickname: params[:nickname].presence
    )

    f(:success)
    redirect_to edit_account_profile_path
  end

  sig { void }
  def destroy
    T.must(current_user).credentials.find(params[:id]).destroy!

    f(:success)
    redirect_to edit_account_profile_path
  end

  private

  sig { params(_exception: T.untyped).void }
  def handle_failure(_exception = nil)
    f(:error, type: :alert)
    redirect_to edit_account_profile_path
  end
end
