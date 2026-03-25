# frozen_string_literal: true

class AdminConstraint
  def matches?(request)
    session_id = request.cookie_jar.signed[:session_id]
    return false unless session_id

    current_session = Session.includes(:user).find_by(id: session_id)
    return false unless current_session

    user = current_session.user
    user&.admin?
  end
end
