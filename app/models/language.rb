# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson'
  has_many :module_descriptions, through: :modules, source: :descriptions
  has_many :lesson_descriptions, through: :lessons, source: :descriptions

  belongs_to :current_version, optional: true, class_name: 'Language::Version'
  has_many :current_lessons, through: :current_version, source: :lesson_versions
  has_many :current_modules, through: :current_version, source: :module_versions

  delegate :to_s, to: :current_version
end
