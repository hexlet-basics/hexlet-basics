class AddAdminsRoleToAleksandrLitvinov < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'ssssank@gmail.com')&.update(admin: true)
  end
end
