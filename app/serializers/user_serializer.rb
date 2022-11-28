# frozen_string_literal: true

class UserSerializer
  class << self
    def to_select2(user)
      Jbuilder.new do |person|
        person.id user.id
        person.text user.to_s
      end.attributes!
    end
  end
end
