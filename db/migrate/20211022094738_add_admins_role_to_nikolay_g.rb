class AddAdminsRoleToNikolayG < ActiveRecord::Migration[6.1]
  def change
    User.find_by(email: "feycot@gmail.com")&.update(admin: true)
  end
end
