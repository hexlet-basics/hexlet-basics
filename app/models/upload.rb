# frozen_string_literal: true

class Upload < ApplicationRecord
  include AASM
  extend Enumerize

  enumerize :language_name, in: %i[php javascript python java ruby]

  has_one :language_version, dependent: :destroy, class_name: 'Language::Version'
  has_many :module_versions, dependent: :destroy, class_name: 'Language::Module::Version'
  has_many :lesson_versions, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'

  aasm :state do
    state :queued, initial: true
    state :running
    state :success
    state :fail

    event :run do
      transitions from: :queued, to: :running
    end

    event :succeed do
      transitions from: :running, to: :success
    end

    event :failed do
      transitions from: :running, to: :fail
    end
  end
end
