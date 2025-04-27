# == Schema Information
#
# Table name: survey_answers
#
#  id             :integer          not null, primary key
#  state          :string
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  survey_id      :integer          not null
#  survey_item_id :integer
#  user_id        :integer          not null
#
# Indexes
#
#  index_survey_answers_on_survey_id              (survey_id)
#  index_survey_answers_on_survey_id_and_user_id  (survey_id,user_id) UNIQUE
#  index_survey_answers_on_survey_item_id         (survey_item_id)
#  index_survey_answers_on_user_id                (user_id)
#
# Foreign Keys
#
#  survey_id       (survey_id => surveys.id)
#  survey_item_id  (survey_item_id => survey_items.id)
#  user_id         (user_id => users.id)
#
require "test_helper"

class Survey::AnswerTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
