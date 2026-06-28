# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: user_credentials
#
#  id          :bigint           not null, primary key
#  nickname    :string
#  public_key  :string           not null
#  sign_count  :bigint           default(0), not null
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  external_id :string           not null
#  user_id     :bigint           not null
#
# Indexes
#
#  index_user_credentials_on_external_id  (external_id) UNIQUE
#  index_user_credentials_on_user_id      (user_id)
#
# Foreign Keys
#
#  fk_rails_...  (user_id => users.id)
#
class User::Credential < ApplicationRecord
  belongs_to :user
end
