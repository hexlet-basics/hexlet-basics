class CreateSurveyScenarios < ActiveRecord::Migration[8.0]
  def change
    create_table :survey_scenarios do |t|
      t.references :survey_item, foreign_key: true
      t.string :name
      t.string :state
      t.string :locale

      # t.string :event_name
      # t.string :event_threshold_count

      t.timestamps
    end
  end
end
