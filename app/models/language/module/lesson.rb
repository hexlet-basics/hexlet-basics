# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  include Language::Module::LessonRepository

  has_paper_trail

  belongs_to :upload
  belongs_to :language
  belongs_to :module
  belongs_to :current_exercise, class_name: 'Language::Module::Lesson::Exercise', optional: true

  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Lesson::Description'
  has_many :exercises, dependent: :destroy, class_name: 'Language::Module::Lesson::Exercise'
end
