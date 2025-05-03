class AddFieldsToUsers < ActiveRecord::Migration[8.0]
  def change
    add_column :users, :contact_value, :string, null: true
    add_column :users, :contact_method, :string, null: true
  end
end
