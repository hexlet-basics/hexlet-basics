class RenameReferenesToUpload < ActiveRecord::Migration[6.0]
  def change
    rename_column :language_versions, :language_upload_id, :upload_id
    rename_column :language_module_versions, :language_upload_id, :upload_id
    rename_column :language_module_lesson_versions, :language_upload_id, :upload_id
  end
end
