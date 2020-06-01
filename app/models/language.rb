# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson'
  has_many :module_descriptions, through: :modules, source: :descriptions
  has_many :lesson_descriptions, through: :lessons, source: :descriptions

  belongs_to :current_version, optional: true, class_name: 'Language::Version'

  def to_s
    current_version.name
  end

  def current_lessons
    current_version.lesson_versions
  end

  def current_modules
    current_version.module_versions
  end
end
