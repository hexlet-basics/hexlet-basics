# frozen_string_literal: true

class Language::Version < ApplicationRecord
  include AASM

  has_many :module_versions, dependent: :destroy,
                             foreign_key: :language_version_id,
                             class_name: 'Language::Module::Version',
                             inverse_of: :language_version
  has_many :lesson_versions, dependent: :destroy,
                             foreign_key: :language_version_id,
                             class_name: 'Language::Lesson::Version',
                             inverse_of: :language_version
  has_many :lesson_infos, dependent: :destroy,
                          foreign_key: :language_version_id,
                          class_name: 'Language::Lesson::Version::Info',
                          inverse_of: :language_version
  has_many :module_infos, dependent: :destroy,
                          foreign_key: :language_version_id,
                          class_name: 'Language::Module::Version::Info',
                          inverse_of: :language_version
  has_many :lessons, through: :lesson_versions

  belongs_to :language

  aasm :state do
    state :created, initial: true
    state :building
    state :built
    state :failed

    event :build do
      transitions from: :created, to: :building
    end

    event :mark_as_built do
      transitions from: :building, to: :built
    end

    event :mark_as_failed do
      transitions to: :failed
    end
  end

  def to_s
    name
  end

  def to_json(*_args)
    attrs = attributes.extract! 'id', 'name', 'created_at'
    attrs.to_json
  end

  def image_tag
    return "lv#{id}" if Rails.env.production?

    :latest
  end
end
