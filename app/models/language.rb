# frozen_string_literal: true

class Language < ApplicationRecord
  extend Enumerize

  enumerize :slug, in: %i[php javascript python java html css racket elixir ruby go]

  validates :slug, presence: true

  belongs_to :version, optional: true

  has_many :uploads, dependent: :destroy
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy
  has_many :versions, dependent: :destroy

  has_many :current_module_data, through: :version, source: :module_data
  has_many :current_lesson_data, through: :version, source: :lesson_data
  has_many :current_lessons, through: :version, source: :lesson_versions
  has_many :current_modules, through: :version, source: :module_versions

  delegate :to_s, to: :version
end
