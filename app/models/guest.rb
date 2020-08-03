# frozen_string_literal: true

class Guest
  def guest?
    true
  end

  def admin?
    false
  end

  def finished_lessons_for_language(_)
    []
  end

  def not_finished_lessons_for_language(language)
    language.current_lessons
  end

  def complete_language?(_)
    false
  end
end
