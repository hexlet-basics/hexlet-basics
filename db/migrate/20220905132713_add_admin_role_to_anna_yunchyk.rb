class AddAdminRoleToAnnaYunchyk < ActiveRecord::Migration[7.0]
  def change
    user = User.find_by email: 'hannayunchik1994@gmail.com'
    return unless user

    user.admin = true
    user.save!
  end
end
