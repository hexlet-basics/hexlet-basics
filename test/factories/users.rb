# frozen_string_literal: true

# == Schema Information
#
# Table name: users
#
#  id                   :bigint           not null, primary key
#  first_name           :string(255)
#  last_name            :string(255)
#  email                :string(255)
#  nickname             :string(255)
#  github_uid           :integer
#  updated_at           :datetime         not null
#  facebook_uid         :string(255)
#  password_digest      :string(255)
#  confirmation_token   :string(255)
#  reset_password_token :string(255)
#  state                :string(255)
#  locale               :string(255)
#  email_delivery_state :string(255)
#  admin                :boolean
#  created_at           :datetime         not null
#  help                 :boolean
#
FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "Email-#{n}@email.com" }
    password { 'password' }
  end
end
