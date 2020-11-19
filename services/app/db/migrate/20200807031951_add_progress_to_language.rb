class AddProgressToLanguage < ActiveRecord::Migration[6.0]
  def change
    add_column :languages, :progress, :string
  end
end
