# == Schema Information
#
# Table name: user_survey_pivots
#
#  id                        :bigint           not null, primary key
#  created_at                :datetime         not null
#  updated_at                :datetime         not null
#  coding_experience_item_id :bigint
#  goal_item_id              :bigint
#  study_plan_item_id        :bigint
#  user_id                   :bigint           not null
#
# Indexes
#
#  index_user_survey_pivots_on_coding_experience_item_id  (coding_experience_item_id)
#  index_user_survey_pivots_on_goal_item_id               (goal_item_id)
#  index_user_survey_pivots_on_study_plan_item_id         (study_plan_item_id)
#  index_user_survey_pivots_on_user_id                    (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (coding_experience_item_id => survey_items.id)
#  fk_rails_...  (goal_item_id => survey_items.id)
#  fk_rails_...  (study_plan_item_id => survey_items.id)
#  fk_rails_...  (user_id => users.id)
#
require "test_helper"

class UserSurveyPivotTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
