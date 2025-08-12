class CreateLanguageLessonReviews < ActiveRecord::Migration[8.0]
  def change
    create_table :language_lesson_reviews do |t|
      t.references :language, null: false, foreign_key: true
      t.references :language_lesson, null: false, foreign_key: true
      t.references :language_lesson_version, null: false, foreign_key: true
      t.references :language_lesson_version_info, null: false, foreign_key: true
      t.text :summary, null: false

      t.timestamps
    end
  end
end
