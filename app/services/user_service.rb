# typed: strict
# frozen_string_literal: true

class UserService < ApplicationService
  class << self
    extend T::Sig

    sig { params(user: User, suffix: T.untyped).void }
    def reset_password!(user, suffix)
      UserMailer.with(user:, suffix:).reset_password.deliver_later
    end

    sig { params(id: T.untyped, struct: UserStruct).returns(Typed::Result[User, User]) }
    def update(id, struct)
      user = User.find(id)
      return fail_with(user) unless user.update(struct.attributes)

      success_with(user)
    end
  end
end
