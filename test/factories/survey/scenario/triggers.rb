# == Schema Information
#
# Table name: survey_scenario_triggers
#
#  id                    :bigint           not null, primary key
#  event_name            :string
#  event_threshold_count :integer
#  created_at            :datetime         not null
#  updated_at            :datetime         not null
#  scenario_id           :bigint           not null
#
# Indexes
#
#  index_survey_scenario_triggers_on_event_name_and_scenario_id  (event_name,scenario_id) UNIQUE
#  index_survey_scenario_triggers_on_scenario_id                 (scenario_id)
#
# Foreign Keys
#
#  fk_rails_...  (scenario_id => survey_scenarios.id)
#
FactoryBot.define do
  factory :survey_scenario_trigger, class: "Survey::Scenario::Trigger" do
    scenario { nil }
    event_name { "MyString" }
    event_threshold_count { 1 }
  end
end
