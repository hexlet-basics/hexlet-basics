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
class BookRequest < ApplicationRecord
  include AASM

  belongs_to :user

  enum :state, { requested: "requested", downloaded: "downloaded" }, default: "requested", suffix: true, validate: true
end
