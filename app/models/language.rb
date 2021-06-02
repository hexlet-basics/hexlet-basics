# frozen_string_literal: true

class Language < ApplicationRecord
  extend Enumerize

  enumerize :slug, in: %i[php javascript python java html css racket elixir ruby go]
  enumerize :progress, in: %i[completed in_development draft], default: :draft, scope: true
  enumerize :learn_as, in: %i[first_language second_language], default: :first_language

  validates :slug, presence: true

  belongs_to :current_version, optional: true, class_name: 'Language::Version'

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lesson_versions, through: :current_version, source: :lesson_versions
  has_many :current_module_versions, through: :current_version, source: :module_versions
  has_many :current_lessons, through: :current_version, source: :lessons

  delegate :to_s, to: :current_version
  delegate :to_json, to: :current_version

  def duration
    lessons.size * 15 / 60
  end
end
