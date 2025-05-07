class CreateSurveyScenarioTriggers < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_scenario_triggers do |t|
      t.references :scenario, null: false, foreign_key: { to_table: :survey_scenarios }
      t.string :event_name
      t.integer :event_threshold_count

      t.timestamps
    end

    add_index :survey_scenario_triggers, [ :event_name, :scenario_id ], unique: true
  end
end
