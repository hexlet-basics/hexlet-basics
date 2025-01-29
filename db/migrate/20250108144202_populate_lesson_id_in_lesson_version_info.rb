class PopulateLessonIdInLessonVersionInfo < ActiveRecord::Migration[8.0]
  def change
    execute <<-SQL.squish
      UPDATE language_lesson_version_infos
      SET language_lesson_id = language_lesson_versions.lesson_id
      FROM language_lesson_versions
      WHERE language_lesson_versions.id = language_lesson_version_infos.version_id
    SQL
  end
end
