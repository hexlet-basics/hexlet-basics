class AddYmClientIdToAhoyVisits < ActiveRecord::Migration[8.1]
  def change
    add_column :ahoy_visits, :ym_client_id, :string
  end
end
