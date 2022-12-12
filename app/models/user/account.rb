# frozen_string_literal: true

# == Schema Information
#
# Table name: user_accounts
#
#  id         :bigint           not null, primary key
#  user_id    :bigint           not null
#  provider   :string(255)      not null
#  uid        :string(255)      not null
#  updated_at :datetime         not null
#  created_at :datetime         not null
#
class User::Account < ApplicationRecord
  belongs_to :user
end
