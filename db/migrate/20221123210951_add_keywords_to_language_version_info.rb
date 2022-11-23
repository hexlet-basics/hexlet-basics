class AddKeywordsToLanguageVersionInfo < ActiveRecord::Migration[7.0]
  def change
    add_column :language_version_infos, :keywords, :string
  end
end
