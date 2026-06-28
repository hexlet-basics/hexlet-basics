# typed: strict

class Language::Lesson::Member::MessagePolicy < ApplicationPolicy
  extend T::Sig

  sig { returns(T::Boolean) }
  def create?
    return false if user.guest?
    messages_count = user.assistant_messages
      .user_role
      .where(created_at: Date.current.beginning_of_day..).count
    messages_count <= 7
  end
end
