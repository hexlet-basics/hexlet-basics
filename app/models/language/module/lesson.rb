# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  include Language::Module::LessonRepository

  belongs_to :language
  belongs_to :module
  belongs_to :current_version, class_name: 'Language::Module::Lesson::Version', optional: true

  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Lesson::Description'
  has_many :versions, dependent: :destroy, class_name: 'Language::Module::Lesson::Version'
end
