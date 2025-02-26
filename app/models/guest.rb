# frozen_string_literal: true

require "securerandom"

class Guest
  def id
    SecureRandom.uuid
  end

  def created_at; end

  def email; end

  def first_name
    ""
  end

  def last_name
    ""
  end

  def guest?
    true
  end

  def admin?
    false
  end

  def locale; end

  def language_members
    Language::Member.none
  end

  def finished_lessons_for_language(_)
    Language::Lesson.none
  end

  def carrotquest_hash; end

  # def not_finished_lessons_for_language(language)
  #   language.current_lesson_infos
  # end
end
