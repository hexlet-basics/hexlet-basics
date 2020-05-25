# frozen_string_literal: true

class Language < ApplicationRecord
  has_many :modules, dependent: :destroy
  has_many :module_descriptions, class_name: 'Language::Module::Description', dependent: :destroy
  has_many :lesson_descriptions, class_name: 'Language::Module::Lesson::Description', dependent: :destroy

  def to_s
    name
  end
end
