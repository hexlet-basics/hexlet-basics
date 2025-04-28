# frozen_string_literal: true

# == Schema Information
#
# Table name: user_accounts
#
#  id         :bigint           not null, primary key
#  provider   :string(255)      not null
#  uid        :string(255)      not null
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  user_id    :bigint           not null
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
class User::Account < ApplicationRecord
  belongs_to :user
end
