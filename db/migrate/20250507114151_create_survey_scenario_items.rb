class CreateSurveyScenarioItems < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_scenario_items do |t|
      t.references :survey, null: false, foreign_key: true
      t.references :scenario, null: false, foreign_key: { to_table: :survey_scenarios }

      t.timestamps
    end

    add_index :survey_scenario_items, [ :survey_id, :scenario_id ], unique: true
  end
end
