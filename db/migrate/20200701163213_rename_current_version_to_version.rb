class RenameCurrentVersionToVersion < ActiveRecord::Migration[6.0]
  def change
    remove_column :languages, :current_version_id
    remove_column :language_modules, :current_version_id
    add_reference :languages, :version, index: true
    add_reference :language_modules, :version, index: true
  end
end
