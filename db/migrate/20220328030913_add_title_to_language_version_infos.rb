class AddTitleToLanguageVersionInfos < ActiveRecord::Migration[7.0]
  def change
    add_column :language_version_infos, :title, :string
  end
end
