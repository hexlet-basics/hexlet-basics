# frozen_string_literal: true

class Language::Version < ApplicationRecord
  include AASM

  has_many :module_versions, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Module::Version'
  has_many :lesson_versions, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Lesson::Version'
  has_many :lesson_infos, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Lesson::Version::Info'
  has_many :module_infos, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Module::Version::Info'
  has_many :lessons, through: :lesson_versions

  belongs_to :language

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

  def to_s
    name
  end
end
