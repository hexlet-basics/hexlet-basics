class AddStateToLanguage < ActiveRecord::Migration[6.0]
  def change
    add_column :languages, :state, :string
  end
end
