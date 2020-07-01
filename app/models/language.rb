# frozen_string_literal: true

class Language < ApplicationRecord
  extend Enumerize

  enumerize :slug, in: %i[php javascript python java html css racket elixir ruby go]

  validates :slug, presence: true

  belongs_to :version, optional: true, class_name: 'Language::Version'

  has_many :uploads, dependent: :destroy
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson'
  has_many :module_descriptions, through: :modules, source: :descriptions
  has_many :lesson_descriptions, through: :lessons, source: :descriptions

  has_many :current_lessons, through: :version, source: :lesson_versions
  has_many :current_modules, through: :version, source: :module_versions

  delegate :to_s, to: :version
end
