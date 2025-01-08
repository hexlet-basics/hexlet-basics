# frozen_string_literal: true

module Language::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { active }
    scope :ordered, -> { order("language_lesson_versions.natural_order") }
  end
end
