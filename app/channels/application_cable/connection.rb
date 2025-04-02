module ApplicationCable
  class Connection < ActionCable::Connection::Base
    identified_by :current_user

    def connect
      self.current_user = find_verified_user
    end

    private
      def find_verified_user
        user_id = (cookies.encrypted[COOKIE_STORE_KEY] || {})["user_id"]
        unless user_id
          reject_unauthorized_connection
          return
        end

        verified_user = User.find_by(id: user_id)
        if verified_user
          verified_user
        else
          reject_unauthorized_connection
        end
      end
  end
end
