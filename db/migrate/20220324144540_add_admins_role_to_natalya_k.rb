class AddAdminsRoleToNatalyaK < ActiveRecord::Migration[6.1]
  def change
    User.find_by(email: "natalya.konstantinova@hexlet.io")&.update(admin: true)
  end
end
