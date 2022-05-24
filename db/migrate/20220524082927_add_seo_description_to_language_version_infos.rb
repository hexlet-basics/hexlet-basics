class AddSeoDescriptionToLanguageVersionInfos < ActiveRecord::Migration[7.0]
  def change
    add_column :language_version_infos, :seo_description, :text
  end
end
