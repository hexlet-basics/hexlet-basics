# == Schema Information
#
# Table name: surveys
#
#  id                    :bigint           not null, primary key
#  description           :string
#  locale                :string
#  question              :string
#  slug                  :string
#  state                 :string
#  created_at            :datetime         not null
#  updated_at            :datetime         not null
#  parent_survey_id      :bigint
#  parent_survey_item_id :bigint
#
# Indexes
#
#  index_surveys_on_parent_survey_id       (parent_survey_id)
#  index_surveys_on_parent_survey_item_id  (parent_survey_item_id)
#  index_surveys_on_slug_and_locale        (slug,locale) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (parent_survey_id => surveys.id)
#  fk_rails_...  (parent_survey_item_id => survey_items.id)
#
require "test_helper"

class SurveyTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
