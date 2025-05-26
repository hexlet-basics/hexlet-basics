class AddAhoyVisitToLeads < ActiveRecord::Migration[8.0]
  def change
    add_reference :leads, :ahoy_visit, null: true, foreign_key: true
  end
end
