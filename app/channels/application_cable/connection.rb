# typed: strict

module ApplicationCable
  class Connection < ActionCable::Connection::Base
    extend T::Sig

    identified_by :current_user

    sig { void }
    def connect
      self.current_user = find_verified_user
    end

    private
      sig { returns(T.nilable(User)) }
      def find_verified_user
        Current.session = Session.includes(:user).find_by(id: cookies.signed[:session_id])
        Current.user || reject_unauthorized_connection
      end
  end
end
