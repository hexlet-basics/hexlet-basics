# == Schema Information
#
# Table name: user_survey_pivots
#
#  id                        :bigint           not null, primary key
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  coding_experience_item_id :integer
#  goal_item_id              :integer
#  user_id                   :integer          not null
#
# Indexes
#
#  index_user_survey_pivots_on_coding_experience_item_id  (coding_experience_item_id)
#  index_user_survey_pivots_on_goal_item_id               (goal_item_id)
#  index_user_survey_pivots_on_user_id                    (user_id)
#
FactoryBot.define do
  factory :user_survey_pivot do
    user { nil }
    coding_experience_item { nil }
    goal_item { nil }
  end
end
