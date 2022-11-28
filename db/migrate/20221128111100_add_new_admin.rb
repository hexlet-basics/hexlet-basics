class AddNewAdmin < ActiveRecord::Migration[7.0]
  def change
    User.find_by(email: 'ralder@yandex.ru')&.update(admin: true)
  end
end
