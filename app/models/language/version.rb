# frozen_string_literal: true

class Language::Version < ApplicationRecord
  has_many :module_versions, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Module::Version'
  has_many :lesson_versions, dependent: :destroy, foreign_key: :language_version_id, class_name: 'Language::Module::Lesson::Version'

  belongs_to :language
  belongs_to :upload

  def to_s
    name
  end
end
