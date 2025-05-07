class CreateSurveyScenarioMembers < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_scenario_members do |t|
      t.references :user, null: false, foreign_key: true
      t.references :scenario, null: false, foreign_key: { to_table: :survey_scenarios }
      t.string :state

      t.timestamps
    end
  end
end
