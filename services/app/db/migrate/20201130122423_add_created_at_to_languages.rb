class AddCreatedAtToLanguages < ActiveRecord::Migration[6.0]
  def change
    add_column :languages, :created_at, :datetime, precision: 6, null: true
    Language.update_all('created_at = inserted_at')
    change_column_null :languages, :created_at, false
    remove_column :languages, :inserted_at
  end
end
