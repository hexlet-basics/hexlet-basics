class AddColumnsToLanguageVersion < ActiveRecord::Migration[6.0]
  def change
    add_column :language_versions, :state, :string
    add_column :language_versions, :result, :string
  end
end
