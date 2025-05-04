# frozen_string_literal: true

require "securerandom"

class Guest
  def id
    SecureRandom.uuid
  end

  def created_at; end

  def email; end

  def assistant_messages
    Language::Lesson::Member::Message.none
  end

  def assistant_messages_count
    0
  end

  def should_add_contact_method
    false
  end

  def book_request
    nil
  end

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

  # def not_finished_lessons_for_language(language)
  #   language.current_lesson_infos
  # end
end
