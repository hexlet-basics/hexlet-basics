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
FactoryBot.define do
  factory :survey_scenario_member, class: "Survey::Scenario::Member" do
    user { nil }
    survey_scenario { nil }
    state { "MyString" }
  end
end
