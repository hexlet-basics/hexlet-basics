class AddIndexOnLanguageLessonVersionInfos < ActiveRecord::Migration[8.0]
  disable_ddl_transaction!

  def change
    add_index :language_lesson_version_infos,
              [ :version_id, :locale ],
              algorithm: :concurrently
  end
end
