# typed: true
# frozen_string_literal: true

module UserRepository
  extend ActiveSupport::Concern
  extend T::Helpers
  requires_ancestor { User }

  included do
    T.bind(self, T.class_of(ActiveRecord::Base))
    scope :admin, -> { where(admin: true) }
  end

  def not_finished_lessons_for_language(language)
    language.current_lessons.left_join_lesson_member_and_user(self)
            .merge(Language::Lesson::Member.started_or_nil)
  end

  def finished_lessons_for_language(language)
    lessons.merge(Language::Lesson::Member.finished).where(id: language.current_lessons)
  end
end
