# typed: strict

# Gates mounted engines (e.g. the Flipper UI) behind admin authentication,
# reusing the app's signed-cookie Session auth (see Authentication concern).
# The cookie jar is populated by ActionDispatch::Cookies before routing.
class AdminConstraint
  extend T::Sig

  sig { params(request: T.untyped).returns(T::Boolean) }
  def self.matches?(request)
    session_id = request.cookie_jar.signed[:session_id]
    return false if session_id.blank?

    session = Session.find_by(id: session_id)
    session&.user&.admin? || false
  end
end
