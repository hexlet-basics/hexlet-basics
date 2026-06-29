# typed: false
# frozen_string_literal: true

# == Schema Information
#
# Table name: users
#
#  id                       :bigint           not null, primary key
#  admin                    :boolean
#  assistant_messages_count :integer          default(0)
#  confirmation_token       :string(255)
#  contact_method           :string
#  contact_value            :string
#  email                    :string(255)
#  email_delivery_state     :string(255)
#  facebook_uid             :string(255)
#  first_name               :string(255)
#  github_uid               :integer
#  help                     :boolean
#  last_name                :string(255)
#  locale                   :string(255)
#  nickname                 :string(255)
#  password_digest          :string(255)
#  phone                    :string
#  phone_verified_at        :datetime
#  state                    :string(255)
#  created_at               :datetime         not null
#  updated_at               :datetime         not null
#  webauthn_id              :string
#
# Indexes
#
#  index_users_on_LOWER_email  (lower((email)::text)) UNIQUE
#  index_users_on_email        (email) UNIQUE
#  index_users_on_phone        (phone) UNIQUE WHERE (phone IS NOT NULL)
#  index_users_on_webauthn_id  (webauthn_id) UNIQUE WHERE (webauthn_id IS NOT NULL)
#
FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "Email-#{n}@email.com" }
    password { "password" }
  end
end
