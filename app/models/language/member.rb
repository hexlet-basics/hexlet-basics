# frozen_string_literal: true

class Language::Member < ApplicationRecord
  include AASM

  belongs_to :user
  belongs_to :language

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: %i[started finished], to: :finished
    end
  end

  def check_for_finish!
    return if finished?

    language_lesson_versions = language.current_lesson_versions

    finished_lessons_in_language =
      language_lesson_versions
      .joins('INNER JOIN language_lesson_members ON language_lesson_members.lesson_id = language_lesson_versions.lesson_id')
      .where('language_lesson_members.user_id = ?', user.id)
      .where(language_lesson_members: { state: :finished })

    return unless language_lesson_versions.count == finished_lessons_in_language.count

    finish!
  end
end
