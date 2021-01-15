class AddAdminFieldToUsers < ActiveRecord::Migration[6.0]
  def change
    add_column :users, :admin, :boolean
  end
end
