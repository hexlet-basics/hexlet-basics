# frozen_string_literal: true

class Language::Lesson::Version < ApplicationRecord
  belongs_to :language_version, class_name: 'Language::Version'
  belongs_to :lesson
  belongs_to :language
  belongs_to :module_version, class_name: 'Language::Module::Version'

  has_many :infos, dependent: :destroy

  def next_lesson
    language_version
      .lesson_versions.order(:natural_order)
      .find_by('natural_order > ?', natural_order)&.lesson
  end

  def prev_lesson
    language_version
      .lesson_versions.order(natural_order: :desc)
      .find_by('natural_order < ?', natural_order)&.lesson
  end
end
