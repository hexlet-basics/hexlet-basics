# frozen_string_literal: true

FactoryBot.define do
  factory :user_account do
    uid { 'MyString' }
    provider { 'MyString' }
    user { nil }
  end
end
