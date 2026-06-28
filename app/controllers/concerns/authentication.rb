# typed: strict

module Authentication
  extend ActiveSupport::Concern
  extend T::Sig
  extend T::Helpers
  requires_ancestor { ApplicationController }

  included do
    T.bind(self, T.class_of(ApplicationController))
    before_action :require_authentication
    # helper_method :authenticated?
  end

  class_methods do
    extend T::Sig

    sig { params(options: T.untyped).returns(T.untyped) }
    def allow_unauthenticated_access(**options)
        T.bind(self, T.class_of(ApplicationController))
      skip_before_action :require_authentication, **options
    end
  end

  sig { params(user: T.untyped).returns(T.untyped) }
  def sign_in(user)
    start_new_session_for(user)
    ahoy.authenticate(user)
  end

  sig { returns(T.untyped) }
  def redirect_if_authenticated
    redirect_to root_path if authenticated?
  end

  sig { returns(T.nilable(User)) }
  def current_user
    resume_session

    Current.user
  end

  sig { returns(T.untyped) }
  def authenticate_admin!
    require_authentication
    redirect_to root_path unless current_user&.admin?
  end

  sig { returns(T.untyped) }
  def require_api_auth!
    head :forbidden unless authenticated?
  end

  sig { returns(T.untyped) }
  def require_admin_api_auth!
    head :forbidden unless authenticated? && current_user&.admin?
  end

  private
    sig { returns(T::Boolean) }
    def authenticated?
      !resume_session.nil?
    end

    sig { returns(T.untyped) }
    def require_authentication
      resume_session || request_authentication
    end

    sig { returns(T.untyped) }
    def resume_session
      Current.session ||= find_session_by_cookie || migrate_legacy_session
    end

    sig { returns(T.untyped) }
    def find_session_by_cookie
      Session.find_by(id: cookies.signed[:session_id]) if cookies.signed[:session_id]
    end

    # TODO: удалить после 10 июля — миграция старых сессий из cookie_store в Session
    sig { returns(T.untyped) }
    def migrate_legacy_session
      user_id = session[:user_id]
      return if user_id.blank?

      user = User.find_by(id: user_id)
      return unless user

      auth_session = start_new_session_for(user)
      session.delete(:user_id)
      auth_session
    end

    sig { returns(T.untyped) }
    def request_authentication
      session[:return_to_after_authenticating] = request.url
      redirect_to new_session_path
    end

    sig { returns(T.untyped) }
    def after_authentication_url
      session.delete(:return_to_after_authenticating) || root_url
    end

    sig { params(user: T.untyped).returns(T.untyped) }
    def start_new_session_for(user)
      user.sessions.create!(user_agent: request.user_agent, ip_address: request.remote_ip).tap do |session|
        Current.session = session
        cookies.signed.permanent[:session_id] = { value: session.id, httponly: true, same_site: :lax }
      end
    end

    sig { returns(T.untyped) }
    def terminate_session
      Current.session.destroy
      cookies.delete(:session_id)
    end
end
