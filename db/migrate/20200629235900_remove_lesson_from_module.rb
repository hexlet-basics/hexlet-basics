class RemoveLessonFromModule < ActiveRecord::Migration[6.0]
  def change
    drop_table :language_module_lesson_descriptions
    drop_table :language_module_lesson_versions
    drop_table :language_module_lessons
    drop_table :language_uploads
    drop_table :language_module_descriptions
  end
end
