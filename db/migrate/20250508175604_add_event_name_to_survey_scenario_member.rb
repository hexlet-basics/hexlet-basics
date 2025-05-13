class AddEventNameToSurveyScenarioMember < ActiveRecord::Migration[8.0]
  def change
    add_column :survey_scenario_members, :event_name, :string
  end
end
