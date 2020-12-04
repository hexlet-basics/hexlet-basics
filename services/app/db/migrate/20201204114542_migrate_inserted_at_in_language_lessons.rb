class MigrateInsertedAtInLanguageLessons < ActiveRecord::Migration[6.0]
  def change
    add_column :language_lessons, :created_at, :datetime, precision: 6, null: true
    Language::Lesson.update_all("created_at = inserted_at")
    change_column_null :language_lessons, :created_at, false
    remove_column :language_lessons, :inserted_at
  end
end
