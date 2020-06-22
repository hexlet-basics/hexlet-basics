# frozen_string_literal: true

class Language::Upload < ApplicationRecord
  include AASM

  belongs_to :language

  has_one  :language_version, dependent: :destroy, class_name: 'Language::Version'
  has_many :language_module_versions, dependent: :destroy, class_name: 'Language::Module::Version'
  has_many :language_module_lesson_versions, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'

  aasm :state do
    state :created, initial: true
    state :building
    state :built

    event :build do
      transitions from: :created, to: :building
    end

    event :done do
      transitions from: :building, to: :built
    end
 end
end
