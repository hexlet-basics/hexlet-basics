# frozen_string_literal: true

class Language::Lesson < ApplicationRecord
  belongs_to :language
  belongs_to :module

  has_many :versions, dependent: :destroy
  has_many :members, dependent: :destroy

  has_many :infos, through: :versions, class_name: 'Language::Lesson::Version::Info'
  has_many :current_infos, through: :current_version, class_name: 'Language::Lesson::Version::Info', source: :infos

  def outdated?(version)
    current_version != version
  end

  def next_lesson
    return nil unless current_version

    current_lesson_natural_order = current_version.natural_order

    current_version
      .language_version
      .lesson_versions.order(:natural_order)
      .where('natural_order > ?', current_lesson_natural_order)
      .limit(1)
      .first&.lesson
  end

  def prev_lesson
    return nil unless current_version

    current_lesson_natural_order = current_version.natural_order

    current_version
      .language_version
      .lesson_versions.order(natural_order: :desc)
      .where('natural_order < ?', current_lesson_natural_order)
      .limit(1)
      .first&.lesson
  end
end
