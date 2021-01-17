class AddDifficultyToLanguage < ActiveRecord::Migration[6.1]
  def change
    add_column :languages, :difficulty, :string
  end
end
