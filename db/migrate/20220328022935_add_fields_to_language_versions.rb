class AddFieldsToLanguageVersions < ActiveRecord::Migration[7.0]
  def change
    add_column :language_versions, :learn_as, :string
    add_column :language_versions, :progress, :string
  end
end
