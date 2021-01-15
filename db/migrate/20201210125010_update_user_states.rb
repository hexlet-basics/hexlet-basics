class UpdateUserStates < ActiveRecord::Migration[6.0]
  def change
    User.where(state: nil).update_all(state: :active)
  end
end
