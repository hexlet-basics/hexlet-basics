class AddYmClientIdToLeads < ActiveRecord::Migration[8.0]
  def change
    add_column :leads, :ym_client_id, :string
  end
end
