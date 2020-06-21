# frozen_string_literal: true

class Language::Upload < ApplicationRecord
  include AASM

  belongs_to :language

  has_one  :language_version, dependent: :destroy, class_name: 'Language::Version'
  has_many :language_module_versions, dependent: :destroy, class_name: 'Language::Module::Version'
  has_many :language_module_lesson_versions, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'

  aasm :state do
    state :not_run, initial: true
    state :running
    state :success
    state :failed

    event :run do
      transitions from: :not_run, to: :running
    end

    event :succeed do
      transitions from: :running, to: :success
    end

    event :fail do
      transitions from: :running, to: :failed
    end
  end
end
