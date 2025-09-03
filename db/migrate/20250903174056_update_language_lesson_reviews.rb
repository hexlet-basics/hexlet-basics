class UpdateLanguageLessonReviews < ActiveRecord::Migration[8.0]
  def change
    change_column_null :language_lesson_reviews, :summary, true
  end
end
