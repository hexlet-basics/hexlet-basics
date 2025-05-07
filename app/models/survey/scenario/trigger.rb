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
class Survey::Scenario::Trigger < ApplicationRecord
  include Survey::Scenario::TriggerRepository

  event_registry = ApplicationContainer["event_registry"]

  belongs_to :scenario

  enum :event_name, event_registry.all.index_with(&:to_s), suffix: true

  validates :event_threshold_count, numericality: { greater_than_or_equal_to: 1 }, allow_nil: true
  validates :event_name, uniqueness: { scope: :scenario }, inclusion: { in: ->(_) { event_registry.all } }, allow_nil: true
end
