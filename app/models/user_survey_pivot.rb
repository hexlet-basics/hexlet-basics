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
class UserSurveyPivot < ApplicationRecord
  belongs_to :user
  belongs_to :coding_experience_item, class_name: "Survey::Item", optional: true
  belongs_to :goal_item, class_name: "Survey::Item", optional: true

  def self.ransackable_attributes(auth_object = nil)
    [ "coding_experience_item_id", "created_at", "goal_item_id", "id", "updated_at", "user_id" ]
  end
end
