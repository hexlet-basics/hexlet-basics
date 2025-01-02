# frozen_string_literal: true

# == Schema Information
#
# Table name: reviews
#
#  id          :integer          not null, primary key
#  body        :text
#  first_name  :string
#  last_name   :string
#  locale      :string
#  state       :string
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  language_id :bigint           not null
#  user_id     :bigint           not null
#
# Indexes
#
#  index_reviews_on_language_id  (language_id)
#  index_reviews_on_user_id      (user_id)
#
# Foreign Keys
#
#  language_id  (language_id => languages.id)
#  user_id      (user_id => users.id)
#
require 'test_helper'

class ReviewTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
