# == Schema Information
#
# Table name: survey_scenario_members
#
#  id          :bigint           not null, primary key
#  event_name  :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  scenario_id :bigint           not null
#  user_id     :bigint           not null
#
# Indexes
#
#  index_survey_scenario_members_on_scenario_id  (scenario_id)
#  index_survey_scenario_members_on_user_id      (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (scenario_id => survey_scenarios.id)
#  fk_rails_...  (user_id => users.id)
#
class Survey::Scenario::Member < ApplicationRecord
  belongs_to :user
  belongs_to :scenario, class_name: "Survey::Scenario"

  event_registry = DepsLocator.current.event_registry

  validates :scenario, uniqueness: { scope: :user }

  enum :event_name, event_registry.all.index_with(&:to_s), suffix: true
  enum :state, { started: "started", finished: "finished" }, default: "started"

  def next_survey
    scenario.surveys.order(order: :asc).where.not(id: user.survey_answers_surveys).first
  end
end
