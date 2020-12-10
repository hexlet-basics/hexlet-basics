class ChangeUniqIndexOnUserAccordingToStates < ActiveRecord::Migration[6.0]
  def change
    remove_index :users, :users_email_index
    add_index :users, :email, where: "state != 'removed'", uniq: true
  end
end
