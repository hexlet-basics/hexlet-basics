class RenameEncryptedPassword < ActiveRecord::Migration[6.0]
  def change
    rename_column :users, :encrypted_password, :password_digest
  end
end
