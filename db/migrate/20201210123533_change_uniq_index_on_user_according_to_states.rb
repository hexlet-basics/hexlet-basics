class ChangeUniqIndexOnUserAccordingToStates < ActiveRecord::Migration[6.0]
  def change
    remove_index :users, name: :users_email_index
    add_index :users, :email, where: "state != 'removed'", unique: true
  end
end
