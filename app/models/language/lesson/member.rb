# frozen_string_literal: true

class Language::Lesson::Member < ApplicationRecord
  include AASM

  belongs_to :user
  belongs_to :language
  belongs_to :lesson
  belongs_to :lesson_version, class_name: 'Language::Lesson::Version'

  aasm :state do
    state :started, initial: true
    state :finished

    event :finish do
      transitions from: :started, to: :finished
    end
  end
end
