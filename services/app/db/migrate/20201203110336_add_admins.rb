class AddAdmins < ActiveRecord::Migration[6.0]
  def change
    User.find_by(nickname: "PlugIN73").update(admin: true)
    User.find_by(email: "mokevnin@gmail.com").update(admin: true)

  end
end
