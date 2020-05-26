# frozen_string_literal: true

class Language::Module::Lesson < ApplicationRecord
  include Language::Module::LessonRepository

  belongs_to :language
  belongs_to :module
  has_many :descriptions, dependent: :destroy, class_name: 'Language::Module::Lesson::Description'

  def file_name_for_exercise
    "#{id}.#{language.extension}"
  end
end
