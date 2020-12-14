# frozen_string_literal: true

class Guest
  def id; end

  def created_at; end

  def email; end

  def guest?
    true
  end

  def admin?
    false
  end

  def language_members
    []
  end

  def finished_lessons_for_language(_)
    []
  end

  def not_finished_lessons_for_language(language)
    language.current_lessons
  end
end
