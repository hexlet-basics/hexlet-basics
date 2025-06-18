# frozen_string_literal: true

module Language::LessonRepository
  extend ActiveSupport::Concern

  included do
    scope :web, -> { active }

    scope :left_join_lesson_member_and_user, ->(user) {
      lesson_members_table = Language::Lesson::Member.arel_table
      lessons_table = Language::Lesson.arel_table

      join_condition = lesson_members_table[:lesson_id].eq(lessons_table[:id]).and(lesson_members_table[:user_id].eq(user.id))

      joins(lessons_table.join(lesson_members_table, Arel::Nodes::OuterJoin).on(join_condition).join_sources)
    }

    scope :ordered, -> { order("language_lesson_versions.natural_order") }
  end
end
