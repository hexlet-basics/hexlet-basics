# frozen_string_literal: true

class Language < ApplicationRecord
  has_paper_trail

  has_many :modules, dependent: :destroy, autosave: true
  has_many :lessons, dependent: :destroy, class_name: 'Language::Module::Lesson', autosave: true
  has_many :module_descriptions, through: :modules, source: :descriptions
  has_many :lesson_descriptions, through: :lessons, source: :descriptions

  belongs_to :upload

  def to_s
    name
  end

  def previous_version
    versions.last.reify(has_many: true, mark_for_destruction: true)
  end

  def previous_version!
    versions.last.reify(has_many: true, mark_for_destruction: true).save!
  end
end
