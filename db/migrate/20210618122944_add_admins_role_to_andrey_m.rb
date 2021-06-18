class AddAdminsRoleToAndreyM < ActiveRecord::Migration[6.1]
  def change
    User.find_by(email: "fakesupp@gmail.com")&.update(admin: true)
  end
end
