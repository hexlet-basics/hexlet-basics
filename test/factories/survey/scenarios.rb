# == Schema Information
#
# Table name: survey_scenarios
#
#  id             :bigint           not null, primary key
#  locale         :string
#  name           :string
#  state          :string
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  survey_item_id :bigint
#
# Indexes
#
#  index_survey_scenarios_on_survey_item_id  (survey_item_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_item_id => survey_items.id)
#
FactoryBot.define do
  factory :survey_scenario, class: "Survey::Scenario" do
    name { "MyString" }
  end
end
