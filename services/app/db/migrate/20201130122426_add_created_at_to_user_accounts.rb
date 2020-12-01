class AddCreatedAtToUserAccounts < ActiveRecord::Migration[6.0]
  def change
    add_column :user_accounts, :created_at, :datetime, precision: 6, null: true
    User::Account.update_all('created_at = inserted_at')
    change_column_null :user_accounts, :created_at, false
    remove_column :user_accounts, :inserted_at
  end
end
