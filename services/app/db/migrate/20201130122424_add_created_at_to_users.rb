class AddCreatedAtToUsers < ActiveRecord::Migration[6.0]
  def change
    add_column :users, :created_at, :datetime, precision: 6, null: true
    User.update_all('created_at = inserted_at')
    change_column_null :users, :created_at, false
    remove_column :users, :inserted_at
  end
end
