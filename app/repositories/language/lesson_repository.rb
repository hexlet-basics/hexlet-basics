# frozen_string_literal: true

module Language::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { active }
    scope :left_join_lesson_member_and_user, ->(user) { joins("LEFT OUTER JOIN language_lesson_members ON language_lesson_members.lesson_id = language_lessons.id AND language_lesson_members.user_id = #{user.id}") }
    scope :ordered, -> { order('language_lesson_versions.natural_order') }
  end
end
