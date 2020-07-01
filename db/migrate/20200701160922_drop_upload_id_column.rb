class DropUploadIdColumn < ActiveRecord::Migration[6.0]
  def change
    remove_column :language_module_versions, :upload_id
    remove_column :language_versions, :upload_id
  end
end
