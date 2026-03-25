module ApplicationCable
  class Connection < ActionCable::Connection::Base
    identified_by :current_user

    def connect
      self.current_user = find_verified_user
    end

    private
      def find_verified_user
        session_id = cookies.signed[:session_id]
        unless session_id
          reject_unauthorized_connection
          return
        end

        verified_session = Session.includes(:user).find_by(id: session_id)
        verified_user = verified_session&.user
        if verified_user
          verified_user
        else
          reject_unauthorized_connection
        end
      end
  end
end
