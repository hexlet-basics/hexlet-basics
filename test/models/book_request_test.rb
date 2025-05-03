# == Schema Information
#
# Table name: book_requests
#
#  id         :bigint           not null, primary key
#  state      :string
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  user_id    :bigint           not null
#
# Indexes
#
#  index_book_requests_on_user_id  (user_id) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
require "test_helper"

class BookRequestTest < ActiveSupport::TestCase
  # test "the truth" do
  #   assert true
  # end
end
