# frozen_string_literal: true

class Language::Module::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :module
  belongs_to :language

  has_many :lesson_versions, dependent: :destroy, foreign_key: :module_version_id, class_name: 'Language::Module::Lesson::Version'
end
