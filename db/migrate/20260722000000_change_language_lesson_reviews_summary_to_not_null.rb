# typed: false
# frozen_string_literal: true

class ChangeLanguageLessonReviewsSummaryToNotNull < ActiveRecord::Migration[8.1]
  def up
    change_column_default :language_lesson_reviews, :summary, from: nil, to: ""
    safety_assured do
      execute("UPDATE language_lesson_reviews SET summary = '' WHERE summary IS NULL")
    end
    add_check_constraint :language_lesson_reviews, "summary IS NOT NULL", name: "language_lesson_reviews_summary_null", validate: false
  end

  def down
    remove_check_constraint :language_lesson_reviews, name: "language_lesson_reviews_summary_null"
    change_column_default :language_lesson_reviews, :summary, from: "", to: nil
  end
end
