# frozen_string_literal: true

class AdminConstraint
  def matches?(request)
    return false unless request.session.key?(:user_id)

    user = User.find_by(id: request.session[:user_id])
    user&.admin?
  end
end
