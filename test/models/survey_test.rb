# == Schema Information
#
# Table name: surveys
#
#  id          :integer          not null, primary key
#  description :string
#  locale      :string
#  question    :string
#  slug        :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_surveys_on_slug_and_locale  (slug,locale) UNIQUE
#
require "test_helper"

class SurveyTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
