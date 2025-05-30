# == Schema Information
#
# Table name: survey_scenario_items
#
#  id          :bigint           not null, primary key
#  order       :integer
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  scenario_id :bigint           not null
#  survey_id   :bigint           not null
#
# Indexes
#
#  index_survey_scenario_items_on_scenario_id                (scenario_id)
#  index_survey_scenario_items_on_survey_id                  (survey_id)
#  index_survey_scenario_items_on_survey_id_and_scenario_id  (survey_id,scenario_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (scenario_id => survey_scenarios.id)
#  fk_rails_...  (survey_id => surveys.id)
#
class Survey::Scenario::Item < ApplicationRecord
  belongs_to :survey
  belongs_to :scenario, class_name: "Survey::Scenario"
end
