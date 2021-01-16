class AddUniqLessonSlugIndex < ActiveRecord::Migration[6.0]
  def change
     add_index :language_lessons, %i[language_id slug], unique: true
  end
end
