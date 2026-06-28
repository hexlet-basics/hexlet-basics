# typed: true

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
  class State < T::Enum
    enums do
      Requested = new("requested")
      Downloaded = new("downloaded")
    end
  end

  include AASM

  belongs_to :user

  typed_enum :state, State, default: State::Requested, suffix: true, validate: true
end
