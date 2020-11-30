class AddCreatedAtToUsers < ActiveRecord::Migration[6.0]
  def change
    add_column :users, :created_at, :datetime, precision: 6, null: false
    User.update_all(created_at: :inserted_at)
    remove_column :users, :inserted_at
  end
end
