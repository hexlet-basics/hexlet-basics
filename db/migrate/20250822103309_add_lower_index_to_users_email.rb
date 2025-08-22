class AddLowerIndexToUsersEmail < ActiveRecord::Migration[8.0]
  disable_ddl_transaction!

  def change
    add_index :users,
              "LOWER(email)",
              unique: true,
              algorithm: :concurrently
  end
end
