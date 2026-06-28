# typed: strict
# frozen_string_literal: true

class UserService < ApplicationService
  class << self
    extend T::Sig

    sig { params(user: User, suffix: T.untyped).void }
    def reset_password!(user, suffix)
      return if user.email.blank?

      UserMailer.with(user:, suffix:).reset_password.deliver_later
    end

    sig { params(id: T.untyped, struct: UserStruct).returns(Typed::Result[User, User]) }
    def update(id, struct)
      user = User.find(id)
      return fail_with(user) unless user.update(struct.attributes)

      success_with(user)
    end

    sig { params(user: User, struct: ProfileStruct).returns(Typed::Result[User, User]) }
    def update_profile(user, struct)
      return fail_with(user) unless user.update(struct.attributes)

      success_with(user)
    end

    sig { params(user: User, struct: PasswordStruct).returns(Typed::Result[User, User]) }
    def update_password(user, struct)
      return fail_with(user) unless user.update(password: struct.password)

      success_with(user)
    end

    sig { params(struct: SignUpStruct, locale: String).returns(Typed::Result[User, User]) }
    def sign_up(struct, locale:)
      user = User.new(struct.attributes.merge(locale:))

      begin
        user.save!
        success_with(user)
      rescue ActiveRecord::RecordInvalid, ActiveRecord::RecordNotUnique
        fail_with(user)
      end
    end
  end
end
