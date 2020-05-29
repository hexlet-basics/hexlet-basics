# frozen_string_literal: true

class Language < ApplicationRecord
  has_paper_trail

  has_many :modules, dependent: :destroy
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson'
  has_many :module_descriptions, through: :modules, source: :descriptions
  has_many :lesson_descriptions, through: :lessons, source: :descriptions

  belongs_to :upload

  def to_s
    name
  end
end
