# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  include Language::Module::LessonRepository

  belongs_to :language
  belongs_to :module
  belongs_to :current_version, optional: true, class_name: 'Language::Module::Lesson::Version'
  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Lesson::Description'
end
