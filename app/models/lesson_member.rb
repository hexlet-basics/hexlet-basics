# frozen_string_literal: true

class LessonMember < ApplicationRecord
  include AASM

  belongs_to :user
  belongs_to :lesson, class_name: 'Language::Module::Lesson'
  belongs_to :lesson_version, class_name: 'Language::Module::Lesson::Version'

  aasm :state do
    state :started, initial: true
    state :finished

    event :complete do
      transitions from: :started, to: :finished
    end
  end
end
