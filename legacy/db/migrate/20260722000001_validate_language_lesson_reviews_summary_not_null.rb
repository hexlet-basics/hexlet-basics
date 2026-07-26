# typed: false
# frozen_string_literal: true

class ValidateLanguageLessonReviewsSummaryNotNull < ActiveRecord::Migration[8.1]
  def up
    validate_check_constraint :language_lesson_reviews, name: "language_lesson_reviews_summary_null"
    change_column_null :language_lesson_reviews, :summary, false
    remove_check_constraint :language_lesson_reviews, name: "language_lesson_reviews_summary_null"
  end

  def down
    add_check_constraint :language_lesson_reviews, "summary IS NOT NULL", name: "language_lesson_reviews_summary_null", validate: false
    change_column_null :language_lesson_reviews, :summary, true
  end
end
