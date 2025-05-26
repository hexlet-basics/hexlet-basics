class AddOrderToSurveyScenarioItems < ActiveRecord::Migration[8.0]
  def change
    add_column :survey_scenario_items, :order, :integer
  end
end
