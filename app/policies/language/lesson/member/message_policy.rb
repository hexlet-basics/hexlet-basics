class Language::Lesson::Member::MessagePolicy < ApplicationPolicy
  def create?
    return false if user.guest?
    messages_count = user.assistant_messages.where(created_at: 24.hours.ago..Time.current).count
    messages_count <= 7
  end
end
