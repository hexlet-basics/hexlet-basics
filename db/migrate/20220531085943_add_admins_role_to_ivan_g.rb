class AddAdminsRoleToIvanG < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'dzencot@gmail.com')&.update(admin: true)
  end
end
