class Language::Lesson::Member::MessagePolicy < ApplicationPolicy
  def create?
    return false if user.guest?
    messages_count = user.assistant_messages.where(created_at: Time.current.beginning_of_day..).count
    messages_count <= 7
  end
end
