# frozen_string_literal: true

module AuthConcern
  def sign_in(user)
    session[:user_id] = user.id
    ahoy.authenticate(user)

    event_data = {
      id: user.id,
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name
    }
    event = UserSignedInEvent.new(data: event_data)
    publish_event(event, user)
  end

  def sign_out
    session.delete(:user_id)
    session.clear
  end

  def signed_in?
    !current_user.guest?
  end

  def current_user
    @current_user ||= (session[:user_id].present? ? User.active.find_by(id: session[:user_id]) : nil) || Guest.new
  end

  def authenticate_user!
    return if signed_in?

    redirect_to new_session_path
  end

  def guests_only!
    return if !signed_in?

    redirect_to root_path
  end

  def authenticate_admin!
    redirect_to root_path unless current_user.admin?
  end

  def require_api_auth!
    head :forbidden unless signed_in?
  end

  def require_admin_api_auth!
    head :forbidden unless current_user.admin?
  end
end
