class AddLessonIdToLanguageLessonVersionInfo < ActiveRecord::Migration[8.0]
  def change
    add_reference :language_lesson_version_infos, :language_lesson, null: true, index: false
  end
end
