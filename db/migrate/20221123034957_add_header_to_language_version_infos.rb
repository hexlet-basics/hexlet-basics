class AddHeaderToLanguageVersionInfos < ActiveRecord::Migration[7.0]
  def change
    add_column :language_version_infos, :header, :string
  end
end
