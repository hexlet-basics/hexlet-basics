# typed: strict

class AiMessagePolicy < ApplicationPolicy
  extend T::Sig
  extend T::Generic

  Record = type_member { { fixed: T.untyped } }

  sig { returns(T::Boolean) }
  def create?
    current = user
    return false if current.nil?

    messages_count = current.ai_messages
      .role_user
      .where(created_at: Date.current.beginning_of_day..).count
    messages_count <= 7
  end
end
