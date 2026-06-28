# typed: strict

class Language::Lesson::Member::MessagePolicy < ApplicationPolicy
  extend T::Sig
  extend T::Generic

  Record = type_member { { fixed: T.untyped } }

  sig { returns(T::Boolean) }
  def create?
    current = user
    return false if current.nil?

    messages_count = current.assistant_messages
      .user_role
      .where(created_at: Date.current.beginning_of_day..).count
    messages_count <= 7
  end
end
