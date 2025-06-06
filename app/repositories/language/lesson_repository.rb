# frozen_string_literal: true

module Language::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { active }
    scope :left_join_lesson_member_and_user, ->(user) {
      left_outer_joins(:members)
        .where(language_lesson_members: { user_id: user.id })
    }
    scope :ordered, -> { order("language_lesson_versions.natural_order") }
  end
end
