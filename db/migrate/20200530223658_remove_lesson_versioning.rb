class RemoveLessonVersioning < ActiveRecord::Migration[6.0]
  def change
    drop_table :language_module_lesson_versions
  end
end
