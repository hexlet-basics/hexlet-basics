# == Schema Information
#
# Table name: survey_items
#
#  id         :bigint           not null, primary key
#  order      :integer          not null
#  slug       :string
#  state      :string
#  value      :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  survey_id  :integer          not null
#
# Indexes
#
#  index_survey_items_on_survey_id  (survey_id)
#
# Foreign Keys
#
#  fk_rails_...  (survey_id => surveys.id)
#
require "test_helper"

class Survey::ItemTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
