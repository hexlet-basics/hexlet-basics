# frozen_string_literal: true

module AuthConcern
  def sign_in(user)
    start_new_session_for(user)
    ahoy.authenticate(user)
  end

  def sign_out
    terminate_session
    session.clear
  end

  def signed_in?
    authenticated?
  end

  def current_user
    resume_session

    Current.user || Guest.new
  end

  def authenticate_user!
    require_authentication
  end

  def guests_only!
    return unless authenticated?

    redirect_to root_path
  end

  def authenticate_admin!
    require_authentication
    redirect_to root_path unless current_user.admin?
  end

  def require_api_auth!
    head :forbidden unless authenticated?
  end

  def require_admin_api_auth!
    head :forbidden unless authenticated? && current_user.admin?
  end
end
