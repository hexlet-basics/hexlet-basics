class AddAdminsRoleToYaroslavM < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'sponaproz@gmail.com')&.update(admin: true)
  end
end
