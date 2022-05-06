class AddAdminsRoleToRomanA < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'rmn.shkv@gmail.com')&.update(admin: true)

  end
end
