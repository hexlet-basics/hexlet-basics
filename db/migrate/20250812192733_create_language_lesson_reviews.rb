class CreateLanguageLessonReviews < ActiveRecord::Migration[8.0]
  def change
    create_table :language_lesson_reviews do |t|
      t.references :language, null: false, foreign_key: true
      t.references :lesson, null: false, foreign_key: true
      t.references :lesson_version, null: false, foreign_key: true
      t.references :lesson_info, null: false, foreign_key: true
      t.text :summary

      t.timestamps
    end
  end
end
