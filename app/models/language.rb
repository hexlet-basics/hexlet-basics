# frozen_string_literal: true

class Language < ApplicationRecord
  extend Enumerize

  enumerize :slug, in: %i[php javascript python java html css racket elixir ruby go]

  validates :slug, presence: true

  belongs_to :current_version, optional: true, class_name: "Language::Version"

  has_many :uploads, dependent: :destroy
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy

  has_many :current_module_infos, through: :current_version, source: :module_infos
  has_many :current_lesson_infos, through: :current_version, source: :lesson_infos
  has_many :current_lessons, through: :current_version, source: :lesson_versions
  has_many :current_modules, through: :current_version, source: :module_versions

  delegate :to_s, to: :current_version
end
