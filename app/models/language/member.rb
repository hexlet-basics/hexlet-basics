# frozen_string_literal: true

class Language::Member < ApplicationRecord
  include AASM

  belongs_to :user
  belongs_to :language

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: %i[started finished], to: :finished, guard: :all_lessons_finished?
    end
  end

  def all_lessons_finished?
    not_finished_lessons = user.not_finished_lessons_for_language(language)

    not_finished_lessons.empty?
  end
end
