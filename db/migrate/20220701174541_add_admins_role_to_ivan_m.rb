class AddAdminsRoleToIvanM < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'sgmdlt@protonmail.com')&.update(admin: true)
  end
end
