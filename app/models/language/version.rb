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
  has_many :infos, dependent: :destroy, foreign_key: 'language_version_id', inverse_of: :language_version

  belongs_to :language
  has_one :current_language, class_name: 'Language', foreign_key: 'current_version_id', dependent: :restrict_with_exception, inverse_of: :current_version

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

  def to_hash(*_args)
    attrs = attributes.extract! 'id', 'name', 'created_at'

    language_info = infos.find_by!(locale: I18n.locale)

    attrs.merge({
                  locale: language_info.locale
                })
  end

  # TODO: move to presenter
  def image_tag
    return "lv#{id}" if Rails.env.production?

    :latest
  end

  def locales
    infos.pluck(:locale).uniq
  end
end
