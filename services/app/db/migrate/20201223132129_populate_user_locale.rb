class PopulateUserLocale < ActiveRecord::Migration[6.1]
  def change
    User.where(locale: nil).update_all(locale: :ru)
  end
end
